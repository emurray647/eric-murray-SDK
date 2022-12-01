package lotrsdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	// "github.com/emurray647/eric-murray-SDK/internal"
)

const (
	apiURL = "https://the-one-api.dev/v2/"
)

type Client interface {
	Books(Filter) ([]Book, error)
	ChapterFromBook(book *Book, filter Filter) ([]Chapter, error)
	Movies(Filter) ([]Movie, error)
	QuoteFromMovie(movie *Movie, filter Filter) ([]Quote, error)
	Characters(filter Filter) ([]Character, error)
	QuoteFromCharacter(character *Character, filter Filter) ([]Quote, error)
	Quotes(filter Filter) ([]Quote, error)
}

type client struct {
	token  string
	apiURL string
}

func NewClient(authToken string) Client {
	return client{
		token:  authToken,
		apiURL: apiURL,
	}
}

func (c client) request(endpoint string, filter Filter) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.apiURL, endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	rawQuery, err := filter.GenerateRawQuery()
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = rawQuery

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request %s failed: %w", req.URL, err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// struct to assist in unmarshalling
type unmarshalStruct[T any] struct {
	Docs []T `json:"docs"`
}

func (c client) Books(filter Filter) ([]Book, error) {
	b, err := c.request("book", filter)
	if err != nil {
		return nil, fmt.Errorf("request for books failed: %w", err)
	}

	return retrieve[Book](b)
}

func (c client) ChapterFromBook(book *Book, filter Filter) ([]Chapter, error) {
	b, err := c.request(fmt.Sprintf("/book/%s/chapter", book.ID), filter)
	if err != nil {
		return nil, fmt.Errorf("request for chapters failed: %w", err)
	}

	return retrieve[Chapter](b)
}

func (c client) Movies(filter Filter) ([]Movie, error) {
	b, err := c.request("movie", filter)
	if err != nil {
		return nil, fmt.Errorf("request for movies failed: %w", err)
	}

	return retrieve[Movie](b)
}

func (c client) QuoteFromMovie(movie *Movie, filter Filter) ([]Quote, error) {
	b, err := c.request(fmt.Sprintf("/movie/%s/quote", movie.ID), filter)
	if err != nil {
		return nil, fmt.Errorf("request for quotes failed: %w", err)
	}

	return retrieve[Quote](b)
}

func (c client) Characters(filter Filter) ([]Character, error) {
	b, err := c.request("character", filter)
	if err != nil {
		return nil, fmt.Errorf("request for characters failed: %w", err)
	}

	return retrieve[Character](b)
}

func (c client) QuoteFromCharacter(character *Character, filter Filter) ([]Quote, error) {
	b, err := c.request(fmt.Sprintf("/character/%s/quote", character.ID), filter)
	if err != nil {
		return nil, fmt.Errorf("request for quotes failed: %w", err)
	}

	return retrieve[Quote](b)
}

func (c client) Quotes(filter Filter) ([]Quote, error) {
	b, err := c.request("quote", filter)
	if err != nil {
		return nil, fmt.Errorf("request for quotes failed: %w", err)
	}

	return retrieve[Quote](b)
}

func (c client) Chapters(filter Filter) ([]Chapter, error) {
	b, err := c.request("chapter", filter)
	if err != nil {
		return nil, fmt.Errorf("request for chapters failed: %w", err)
	}

	return retrieve[Chapter](b)
}

func retrieve[T any](b []byte) ([]T, error) {
	data := unmarshalStruct[T]{}
	err := json.Unmarshal(b, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal bytes: %w", err)
	}

	return data.Docs, nil
}

// func (c client) Quotes()
