package lotrsdk

import (
	"encoding/json"
	"fmt"
)

// this file contains the structs representing the JSON response we get from the API

type Book struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

type Movie struct {
	ID                         string  `json:"_id"`
	Name                       string  `json:"name"`
	RuntimeInMinutes           int     `json:"runtimeInMinutes"`
	BudgetInMillions           float32 `json:"budgetInMillions"`
	BoxOfficeRevenueInMillions float32 `json:"boxOfficeRevenueInMillions"`
	AcademyAwardNominations    int     `json:"academyAwardNominations"`
	AcademyAwardWins           int     `json:"academyAwardWins"`
	RottenTomatoesScore        float32 `json:"rottenTomatoesScore"`
}

type Character struct {
	ID      string `json:"_id"`
	Birth   string `json:"birth"`
	Death   string `json:"death"`
	Hair    string `json:"hair"`
	Realm   string `json:"realm"`
	Height  string `json:"height"`
	Spouse  string `json:"spouse"`
	Gender  string `json:"gender"`
	Name    string `json:"name"`
	Race    string `json:"race"`
	WikiURL string `json:"wikiUrl"`
}

type Quote struct {
	ID        string `json:"_id"`
	Dialog    string `json:"dialog"`
	Movie     string `json:"movie"`
	Character string `json:"character"`
}

type Chapter struct {
	ID          string `json:"_id"`
	ChapterName string `json:"chapterName"`
	Book        string `json:"book"`
}

// Status is kept separate from the rest of the structs as a user "probably" doesn't want to deal
// with it, but we can still provide it to them in case they do
type Status struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Page   int `json:"page"`
	Pages  int `json:"pages"`
}

// struct to assist in unmarshalling
type unmarshalStruct[T any] struct {
	Docs []T `json:"docs"`
	Status
}

// unmarshalJSON is a helper function that reads in a byte slice and puts the data into the approriate struct
//   T - the type we are reading
//   b - the byte array from which to read
func unmarshalJSON[T any](b []byte) ([]T, Status, error) {
	data := unmarshalStruct[T]{}
	err := json.Unmarshal(b, &data)
	if err != nil {
		return nil, Status{}, fmt.Errorf("failed to unmarshal bytes: %w", err)
	}

	return data.Docs, data.Status, nil
}
