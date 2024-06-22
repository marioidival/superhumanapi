package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/marioidival/superhuman-api/internal/db"
	"github.com/marioidival/superhuman-api/pkg/database"
)

type Repo struct {
	dbc *database.Client
	q   *db.Queries
}

func NewPersonRepo(dbc *database.Client) Repo {
	return Repo{dbc: dbc, q: db.New(dbc)}
}

func (r Repo) GetPersonByEmail(ctx context.Context, email string) (db.Person, error) {
	person, err := r.q.PersonByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Person{}, database.ErrNotFoundIdentifier
		}
		return db.Person{}, err
	}

	return db.Person{
		ID:       person.ID,
		FullName: person.FullName,
		Avatar:   person.Avatar,
		Location: person.Location,
	}, nil
}

func (r Repo) CreatePerson(ctx context.Context, person db.CreatePersonParams) (db.Person, error) {
	newPerson, err := r.q.CreatePerson(ctx, person)
	if err != nil {
		return db.Person{}, err
	}
	return db.Person{
		ID:       newPerson.ID,
		FullName: newPerson.FullName,
		Email: newPerson.Email,
		Avatar:   newPerson.Avatar,
		Location: newPerson.Location,
	}, nil
}

func (r Repo) UpdatePerson(ctx context.Context, personID uuid.UUID) (int32, error) {
	count, err := r.q.UpdatePersonRanking(ctx, personID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r Repo) PersonPopularity(ctx context.Context) ([]db.PopularPersonsRow, error) {
	persons, err := r.q.PopularPersons(ctx)
	return persons, err
}
