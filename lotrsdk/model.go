package lotrsdk

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

// {
// 	"_id": "5cd95395de30eff6ebccde56",
// 	"name": "The Lord of the Rings Series",
// 	"runtimeInMinutes": 558,
// 	"budgetInMillions": 281,
// 	"boxOfficeRevenueInMillions": 2917,
// 	"academyAwardNominations": 30,
// 	"academyAwardWins": 17,
// 	"rottenTomatoesScore": 94
//   },

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

// "_id": "5cdbdecb6dc0baeae48cfac1",
// "death": "February 293019",
// "birth": " ,3019",
// "hair": "Dark",
// "realm": "NaN",
// "height": "6'1 (film)",
// "spouse": "NaN",
// "gender": "Male",
// "name": "Lugdush",
// "race": "Uruk-hai"

type Quote struct {
	ID        string `json:"_id"`
	Dialog    string `json:"dialog"`
	Movie     string `json:"movie"`
	Character string `json:"character"`
}

// "_id": "5cd96e05de30eff6ebccebd0",
// "dialog": "Get the wounded on horses.The wolves of lsengard will return.Leave the dead.",
// "movie": "5cd95395de30eff6ebccde5b",
// "character": "5cd99d4bde30eff6ebccfe19",
// "id": "5cd96e05de30eff6ebccebd0"

type Chapter struct {
	ID          string `json:"_id"`
	ChapterName string `json:"chapterName"`
	Book        string `json:"book"`
}

// "_id": "6091b6d6d58360f988133bc8",
// "chapterName": "The Grey Havens",
// "book": "5cf58080b53e011a64671584"

// 24 hex characters

// {
// 	"_id": "5cf5805fb53e011a64671582",
// 	"name": "The Fellowship Of The Ring"
//   },
