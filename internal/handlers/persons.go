package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/marioidival/superhuman-api/internal/db"
	"github.com/marioidival/superhuman-api/internal/repository"
	"github.com/marioidival/superhuman-api/pkg/clearbit"
	"github.com/marioidival/superhuman-api/pkg/database"
)

type Server struct {
	personRepo     repository.Repo
	clearbitClient *clearbit.Client
}

type personResponse struct {
	FullName string `json:"full_name"`
	Avatar   string `json:"avatar"`
	Location string `json:"location"`
	Count    int32  `json:"count"`
}


func NewServer(dbc *database.Client, clearbitApiKey string) *Server {
	return &Server{
		personRepo:     repository.NewPersonRepo(dbc),
		clearbitClient: clearbit.NewClient(clearbitApiKey),
	}
}

func (s *Server) EmailLookupHandler(c echo.Context) error {
	ctx := c.Request().Context()
	log.Println("email lookup")

	email := c.QueryParam("email")
	if email == "" {
		log.Println("email query param is required")
		return echo.NewHTTPError(http.StatusBadRequest, "email query param is required")
	}

	person, err := s.personRepo.GetPersonByEmail(ctx, email)
	if err == nil {
		log.Println("person found on database")
		count, errUpdate := s.personRepo.UpdatePerson(ctx, person.ID)
		if errUpdate != nil {
			log.Printf("failed to update search count of person: %s", errUpdate.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, errUpdate.Error())
		}

		return c.JSON(http.StatusOK, personResponse{
			FullName: person.FullName,
			Avatar:   person.Avatar,
			Location: person.Location,
			Count:    count,
		})
	}

	log.Println("getting person from clearbit")

	profile, err := s.clearbitClient.GetPerson(email)
	if err != nil {
		var statusCode int = http.StatusInternalServerError

		if errors.Is(err, clearbit.ErrNotFoundPerson) {
			statusCode = http.StatusNotFound
		}
		log.Printf("failed to get person profile data: %s", err.Error())
		return echo.NewHTTPError(statusCode, err.Error())
	}

	profileByte, err := json.Marshal(profile)
	if err != nil {
		log.Println("failed to marshal person profile data")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	newPerson, err := s.personRepo.CreatePerson(ctx, db.CreatePersonParams{
		FullName: profile.Name.FullName,
		Avatar:   profile.Avatar,
		Email:    profile.Email,
		Location: profile.Location,
		Metadata: profileByte,
	})
	if err != nil {
		log.Printf("failed to save person profile data: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	count, errUpdate := s.personRepo.UpdatePerson(ctx, newPerson.ID)
	if errUpdate != nil {
		log.Printf("failed to update search count of person: %s", errUpdate.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, errUpdate.Error())
	}

	return c.JSON(http.StatusOK, personResponse{
		FullName: newPerson.FullName,
		Avatar:   newPerson.Avatar,
		Location: newPerson.Location,
		Count:    count,
	})
}

func (s *Server) PopularityHandler(c echo.Context) error {
	ctx := c.Request().Context()
	log.Println("popularity")
	persons, err := s.personRepo.PersonPopularity(ctx)
	if err != nil {
		return err
	}
	response := make([]personResponse, 0)
	for _, p := range persons {
		response = append(response, personResponse{
			FullName: p.FullName,
			Avatar: p.Avatar,
			Location: p.Location,
			Count: p.SearchCount,
		})
	}
	return c.JSON(200, response)
}
