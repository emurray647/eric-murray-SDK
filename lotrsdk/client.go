package lotrsdk

import (
	"fmt"
	"io"
	"net/http"
)

const (
	apiURL = "https://the-one-api.dev/v2"
)

// Client exposes several methods to access the-one-api
// each method return a slice of the type it was searching for, the response status, and an error
type Client interface {
	// Books retrieves all the LOTR books
	//   filter - any number of Filter objects
	Books(filter ...Filter) ([]Book, Status, error)

	// ChapterFromBook retrieves all the chapters from the provided book
	//   book - the book from which to get the chapters
	//   filter - any number of filters objects
	ChapterFromBook(book *Book, filter ...Filter) ([]Chapter, Status, error)

	// Movies retrieves all the LOTR movies
	//   filter - any number of Filter objects
	Movies(filter ...Filter) ([]Movie, Status, error)

	// QuoteFromMovie retrieves all the quotes from the provided movie
	//   move - the movie from which to get the quotes
	//   filter - any number of Filter objects
	QuoteFromMovie(movie *Movie, filter ...Filter) ([]Quote, Status, error)

	// Characters retrieves all the characters from the LOTR
	//   filter - any number of Filter objects
	Characters(filter ...Filter) ([]Character, Status, error)

	// QuoteFromCharacter retrieves all the quotes from a character
	//   character - the character who spoke the quote
	//   filter - any number of Filter objects
	QuoteFromCharacter(character *Character, filter ...Filter) ([]Quote, Status, error)

	// Quotes retrieves all the LOTR quotes
	//   filter - any number of Filter objects
	Quotes(filter ...Filter) ([]Quote, Status, error)

	// Chapters retrieves all the LOTR chapters
	//   filter - any number of Filter objects
	Chapters(filter ...Filter) ([]Chapter, Status, error)
}

// client is a Client implementation
type client struct {
	token  string
	apiURL string
}

// NewClient creates a new Client
// authToken - the-one-api authentication token
func NewClient(authToken string) Client {
	return client{
		token:  authToken,
		apiURL: apiURL,
	}
}

// helper function to perform the request
// returns a byte array of the response JSON
func (c client) doRequest(endpoint string, filter ...Filter) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.apiURL, endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	rawQuery, err := MergeFilters(filter...).GenerateRawQuery()
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = rawQuery

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request %s failed: %w", req.URL, err)
	} else if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request %s failed with status code %d:%s", req.URL, resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c client) Books(filter ...Filter) ([]Book, Status, error) {
	b, err := c.doRequest("/book", filter...)
	if err != nil {
		return nil, Status{}, fmt.Errorf("request for books failed: %w", err)
	}

	return unmarshalJSON[Book](b)
}

func (c client) ChapterFromBook(book *Book, filter ...Filter) ([]Chapter, Status, error) {
	b, err := c.doRequest(fmt.Sprintf("/book/%s/chapter", book.ID), filter...)
	if err != nil {
		return nil, Status{}, fmt.Errorf("request for chapters failed: %w", err)
	}

	return unmarshalJSON[Chapter](b)
}

func (c client) Movies(filter ...Filter) ([]Movie, Status, error) {
	b, err := c.doRequest("/movie", filter...)
	if err != nil {
		return nil, Status{}, fmt.Errorf("request for movies failed: %w", err)
	}

	return unmarshalJSON[Movie](b)
}

func (c client) QuoteFromMovie(movie *Movie, filter ...Filter) ([]Quote, Status, error) {
	b, err := c.doRequest(fmt.Sprintf("/movie/%s/quote", movie.ID), filter...)
	if err != nil {
		return nil, Status{}, fmt.Errorf("request for quotes failed: %w", err)
	}

	return unmarshalJSON[Quote](b)
}

func (c client) Characters(filter ...Filter) ([]Character, Status, error) {
	b, err := c.doRequest("/character", filter...)
	if err != nil {
		return nil, Status{}, fmt.Errorf("request for characters failed: %w", err)
	}

	return unmarshalJSON[Character](b)
}

func (c client) QuoteFromCharacter(character *Character, filter ...Filter) ([]Quote, Status, error) {
	b, err := c.doRequest(fmt.Sprintf("/character/%s/quote", character.ID), filter...)
	if err != nil {
		return nil, Status{}, fmt.Errorf("request for quotes failed: %w", err)
	}

	return unmarshalJSON[Quote](b)
}

func (c client) Quotes(filter ...Filter) ([]Quote, Status, error) {
	b, err := c.doRequest("/quote", filter...)
	if err != nil {
		return nil, Status{}, fmt.Errorf("request for quotes failed: %w", err)
	}

	return unmarshalJSON[Quote](b)
}

func (c client) Chapters(filter ...Filter) ([]Chapter, Status, error) {
	b, err := c.doRequest("/chapter", filter...)
	if err != nil {
		return nil, Status{}, fmt.Errorf("request for chapters failed: %w", err)
	}

	return unmarshalJSON[Chapter](b)
}
