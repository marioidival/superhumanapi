package clearbit

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

var ErrNotFoundPerson = errors.New("person not found")

type Person struct {
	ID   string `json:"id"`
	Name struct {
		FullName   string `json:"fullName"`
		GivenName  string `json:"givenName"`
		FamilyName string `json:"familyName"`
	} `json:"name"`
	Email     string `json:"email"`
	Location  string `json:"location"`
	TimeZone  string `json:"timeZone"`
	UtcOffset int    `json:"utcOffset"`
	Geo       struct {
		City        string  `json:"city"`
		State       string  `json:"state"`
		StateCode   string  `json:"stateCode"`
		Country     string  `json:"country"`
		CountryCode string  `json:"countryCode"`
		Lat         float64 `json:"lat"`
		Lng         float64 `json:"lng"`
	} `json:"geo"`
	Bio        interface{} `json:"bio"`
	Site       string      `json:"site"`
	Avatar     string      `json:"avatar"`
	Employment struct {
		Domain    string `json:"domain"`
		Name      string `json:"name"`
		Title     string `json:"title"`
		Role      string `json:"role"`
		Seniority string `json:"seniority"`
	} `json:"employment"`
	Facebook struct {
		Handle string `json:"handle"`
	} `json:"facebook"`
	Github struct {
		Handle    string `json:"handle"`
		ID        int    `json:"id"`
		Avatar    string `json:"avatar"`
		Company   string `json:"company"`
		Blog      string `json:"blog"`
		Followers int    `json:"followers"`
		Following int    `json:"following"`
	} `json:"github"`
	Twitter struct {
		Handle    interface{} `json:"handle"`
		ID        interface{} `json:"id"`
		Bio       interface{} `json:"bio"`
		Followers interface{} `json:"followers"`
		Following interface{} `json:"following"`
		Statuses  interface{} `json:"statuses"`
		Favorites interface{} `json:"favorites"`
		Location  interface{} `json:"location"`
		Site      interface{} `json:"site"`
		Avatar    interface{} `json:"avatar"`
	} `json:"twitter"`
	Linkedin struct {
		Handle string `json:"handle"`
	} `json:"linkedin"`
	Googleplus struct {
		Handle string `json:"handle"`
	} `json:"googleplus"`
	Gravatar struct {
		Handle  string        `json:"handle"`
		Urls    []interface{} `json:"urls"`
		Avatar  string        `json:"avatar"`
		Avatars []struct {
			URL  string `json:"url"`
			Type string `json:"type"`
		} `json:"avatars"`
	} `json:"gravatar"`
	Fuzzy         bool      `json:"fuzzy"`
	EmailProvider bool      `json:"emailProvider"`
	IndexedAt     time.Time `json:"indexedAt"`
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

type Client struct {
	apiKey string
}

func (c *Client) GetPerson(email string) (Person, error) {
	req, err := http.NewRequest("GET", "https://person.clearbit.com/v2/people/find?email="+email, nil)
	if err != nil {
		return Person{}, err
	}

	req.Header.Add("Authorization", "Bearer "+c.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Person{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return Person{}, ErrNotFoundPerson
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Person{}, err
	}

	var person Person
	err = json.Unmarshal(body, &person)
	if err != nil {
		return Person{}, err
	}

	return person, nil
}
