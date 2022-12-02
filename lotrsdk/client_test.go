package lotrsdk

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestOneRingClient() (Client, *[]*http.Request) {

	requests := make([]*http.Request, 0)

	// ts := httptest.NewServer(handlerFunc)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println("received something")
		// fmt.Println(r.RequestURI)
		requests = append(requests, r)
	}))

	client := client{
		token:  "fake-token",
		apiURL: ts.URL,
	}

	return client, &requests
}

func assertQueryContains(t *testing.T, r *http.Request, str string) bool {
	// would like to test the escaped string; unfortunately can't find if Go
	// is standarized on this (ex %20 vs +)
	param, _ := url.QueryUnescape(r.URL.RawQuery)
	return assert.Truef(t, strings.Contains(param, str),
		"expected to contain %s; got %s", str, param)
}

func TestBook(t *testing.T) {
	client, requests := newTestOneRingClient()
	client.Books()

	assert.Equal(t, len(*requests), 1)
	assert.Equal(t, (*requests)[0].URL.Path, "/book")
	assert.Equal(t, len((*requests)[0].URL.Query()), 0)
}

func TestBookChapter(t *testing.T) {
	client, requests := newTestOneRingClient()
	book := Book{
		ID:   "47",
		Name: "The Two Towers",
	}

	client.ChapterFromBook(&book)

	assert.Equal(t, len(*requests), 1)
	assert.Equal(t, (*requests)[0].URL.Path, "/book/47/chapter")
	assert.Equal(t, len((*requests)[0].URL.Query()), 0)
}

func TestMovie(t *testing.T) {
	client, requests := newTestOneRingClient()
	client.Movies()

	assert.Equal(t, len(*requests), 1)
	assert.Equal(t, (*requests)[0].URL.Path, "/movie")
	assert.Equal(t, len((*requests)[0].URL.Query()), 0)
}

func TestMovieQuote(t *testing.T) {
	client, requests := newTestOneRingClient()
	movie := Movie{
		ID:               "501",
		Name:             "The Two Towers",
		AcademyAwardWins: 12,
	}

	client.QuoteFromMovie(&movie)

	assert.Equal(t, len(*requests), 1)
	assert.Equal(t, (*requests)[0].URL.Path, "/movie/501/quote")
	assert.Equal(t, len((*requests)[0].URL.Query()), 0)
}

func TestCharacter(t *testing.T) {
	client, requests := newTestOneRingClient()
	client.Characters()

	assert.Equal(t, len(*requests), 1)
	assert.Equal(t, (*requests)[0].URL.Path, "/character")
	assert.Equal(t, len((*requests)[0].URL.Query()), 0)
}

func TestCharacterQuote(t *testing.T) {
	client, requests := newTestOneRingClient()
	character := Character{
		ID:   "21F3C",
		Name: "Tom Bombadil",
		Hair: "Brown",
	}

	client.QuoteFromCharacter(&character)

	assert.Equal(t, len(*requests), 1)
	assert.Equal(t, (*requests)[0].URL.Path, "/character/21F3C/quote")
	assert.Equal(t, len((*requests)[0].URL.Query()), 0)
}

func TestQuote(t *testing.T) {
	client, requests := newTestOneRingClient()
	client.Quotes()

	assert.Equal(t, len(*requests), 1)
	assert.Equal(t, (*requests)[0].URL.Path, "/quote")
	assert.Equal(t, len((*requests)[0].URL.Query()), 0)
}

func TestChapter(t *testing.T) {
	client, requests := newTestOneRingClient()
	client.Chapters()

	assert.Equal(t, len(*requests), 1)
	assert.Equal(t, (*requests)[0].URL.Path, "/chapter")
	assert.Equal(t, len((*requests)[0].URL.Query()), 0)
}

func TestMatchFilters(t *testing.T) {
	client, requests := newTestOneRingClient()
	client.Characters(BinaryFilter("name", FilterCompareEqual, "Elrond"))
	client.Characters(BinaryFilter("name", FilterCompareNotEqual, "Glorfindel"))

	assert.Equal(t, len(*requests), 2)
	assert.Equal(t, (*requests)[0].URL.Path, "/character")
	assert.Equal(t, len((*requests)[0].URL.Query()), 1)
	assertQueryContains(t, (*requests)[0], "name=Elrond")
	assert.Equal(t, len((*requests)[1].URL.Query()), 1)
	assertQueryContains(t, (*requests)[1], "name!=Glorfindel")
}

func TestIncludeFilters(t *testing.T) {
	client, requests := newTestOneRingClient()
	client.Movies(BinaryFilter("name", FilterCompareEqual, "The Two Towers", "The Battle of the Five Armies"))
	client.Movies(BinaryFilter("name", FilterCompareNotEqual, "The Fellowship of the Ring", "The Hobbit"))

	assert.Equal(t, len(*requests), 2)
	assert.Equal(t, (*requests)[0].URL.Path, "/movie")
	assert.Equal(t, len((*requests)[0].URL.Query()), 1)
	assertQueryContains(t, (*requests)[0], `name=The Two Towers,The Battle of the Five Armies`)
	assert.Equal(t, len((*requests)[1].URL.Query()), 1)
	assertQueryContains(t, (*requests)[1], `name!=The Fellowship of the Ring,The Hobbit`)
}

func TestExists(t *testing.T) {
	client, requests := newTestOneRingClient()
	client.Characters(MergeFilters(ExistFilter("wikiURL"), NotExistFilter("hair")))

	assert.Equal(t, len(*requests), 1)
	assert.Equal(t, (*requests)[0].URL.Path, "/character")
	assertQueryContains(t, (*requests)[0], "wikiURL")
	assertQueryContains(t, (*requests)[0], "!hair")
}

func TestRegex(t *testing.T) {
	client, requests := newTestOneRingClient()
	client.Characters(BinaryFilter("name", FilterCompareEqual, "/foot/i"))

	assert.Equal(t, len(*requests), 1)
	assert.Equal(t, (*requests)[0].URL.Path, "/character")
	assertQueryContains(t, (*requests)[0], "name=/foot/i")
}

func TestInequalities(t *testing.T) {
	client, requests := newTestOneRingClient()
	client.Movies(BinaryFilter("budgetInMillions", FilterCompareLessThan, "100"))
	client.Movies(BinaryFilter("academyAwardWins", FilterCompareLessThanOrEqual, "0"))
	client.Movies(BinaryFilter("runtimeInMinutes", FilterCompareGreaterThan, "160"))
	client.Movies(BinaryFilter("budgetInMillions", FilterCompareGreaterThanOrEqual, "200"))

	assert.Equal(t, len(*requests), 4)
	assert.Equal(t, (*requests)[0].URL.Path, "/movie")
	assertQueryContains(t, (*requests)[0], "budgetInMillions<100")
	assert.Equal(t, (*requests)[1].URL.Path, "/movie")
	assertQueryContains(t, (*requests)[1], "academyAwardWins<=0")
	assert.Equal(t, (*requests)[2].URL.Path, "/movie")
	assertQueryContains(t, (*requests)[2], "runtimeInMinutes>160")
	assert.Equal(t, (*requests)[3].URL.Path, "/movie")
	assertQueryContains(t, (*requests)[3], "budgetInMillions>=200")
}

func TestSorting(t *testing.T) {
	client, requests := newTestOneRingClient()
	client.Characters(Sort("name", SortOrderAscending))
	client.Characters(Sort("hair", SortOrderDescending))

	assert.Equal(t, len(*requests), 2)
	assert.Equal(t, (*requests)[0].URL.Path, "/character")
	assertQueryContains(t, (*requests)[0], "sort=name:asc")
	assert.Equal(t, (*requests)[1].URL.Path, "/character")
	assertQueryContains(t, (*requests)[1], "sort=hair:desc")
}

func TestLimit(t *testing.T) {
	client, requests := newTestOneRingClient()
	client.Books(Limit(2))

	assert.Equal(t, len(*requests), 1)
	assert.Equal(t, (*requests)[0].URL.Path, "/book")
	assertQueryContains(t, (*requests)[0], "limit=2")
}

func TestPage(t *testing.T) {
	client, requests := newTestOneRingClient()
	client.Books(Page(7))

	assert.Equal(t, len(*requests), 1)
	assert.Equal(t, (*requests)[0].URL.Path, "/book")
	assertQueryContains(t, (*requests)[0], "page=7")
}

func TestOffset(t *testing.T) {
	client, requests := newTestOneRingClient()
	client.Books(Offset(31))

	assert.Equal(t, len(*requests), 1)
	assert.Equal(t, (*requests)[0].URL.Path, "/book")
	assertQueryContains(t, (*requests)[0], "offset=31")
}

func newTestClientWithMockServer(data string) Client {

	requests := make([]*http.Request, 0)

	// ts := httptest.NewServer(handlerFunc)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println("received something")
		// fmt.Println(r.RequestURI)
		requests = append(requests, r)
		w.Write([]byte(data))
	}))

	client := client{
		token:  "fake-token",
		apiURL: ts.URL,
	}

	return client
}

func TestUnmarshalBook(t *testing.T) {
	data := `{"docs":[{"_id":"5cf5805fb53e011a64671582","name":"The Fellowship Of The Ring"},{"_id":"5cf58077b53e011a64671583","name":"The Two Towers"},{"_id":"5cf58080b53e011a64671584","name":"The Return Of The King"}],"total":3,"limit":1000,"offset":0,"page":1,"pages":1}`
	client := newTestClientWithMockServer(data)
	books, status, err := client.Books()

	assert.Nil(t, err)
	assert.Equal(t, len(books), 3)
	assert.Equal(t, books[0].ID, "5cf5805fb53e011a64671582")
	assert.Equal(t, books[0].Name, "The Fellowship Of The Ring")
	assert.Equal(t, books[1].ID, "5cf58077b53e011a64671583")
	assert.Equal(t, books[1].Name, "The Two Towers")
	assert.Equal(t, books[2].ID, "5cf58080b53e011a64671584")
	assert.Equal(t, books[2].Name, "The Return Of The King")

	assert.Equal(t, status.Total, 3)
	assert.Equal(t, status.Limit, 1000)
	assert.Equal(t, status.Offset, 0)
	assert.Equal(t, status.Page, 1)
	assert.Equal(t, status.Pages, 1)
}

func TestUnmarshalMovie(t *testing.T) {
	data := `{"docs":[{"_id":"5cd95395de30eff6ebccde56","name":"The Lord of the Rings Series","runtimeInMinutes":558,"budgetInMillions":281,"boxOfficeRevenueInMillions":2917,"academyAwardNominations":30,"academyAwardWins":17,"rottenTomatoesScore":94}],"total":8,"limit":1,"offset":0,"page":1,"pages":8}`
	client := newTestClientWithMockServer(data)
	movies, _, err := client.Movies()

	assert.Nil(t, err)
	assert.Equal(t, len(movies), 1)
	assert.Equal(t, movies[0].ID, "5cd95395de30eff6ebccde56")
	assert.Equal(t, movies[0].RuntimeInMinutes, 558)
}

func TestUnmarshalCharacter(t *testing.T) {
	data := `{"docs":[{"_id":"5cd99d4bde30eff6ebccfbbe","height":"","race":"Human","gender":"Female","birth":"","spouse":"Belemir","death":"","realm":"","hair":"","name":"Adanel","wikiUrl":"http://lotr.wikia.com//wiki/Adanel"}],"total":933,"limit":1,"offset":0,"page":1,"pages":933}`
	client := newTestClientWithMockServer(data)
	characters, _, err := client.Characters()

	assert.Nil(t, err)
	assert.Equal(t, len(characters), 1)
	assert.Equal(t, characters[0].ID, "5cd99d4bde30eff6ebccfbbe")
	assert.Equal(t, characters[0].Race, "Human")
}

func TestUnmarshalQuote(t *testing.T) {
	data := `{"docs":[{"_id":"5cd96e05de30eff6ebcce7e9","dialog":"Deagol!","movie":"5cd95395de30eff6ebccde5d","character":"5cd99d4bde30eff6ebccfe9e","id":"5cd96e05de30eff6ebcce7e9"}],"total":2390,"limit":1,"offset":0,"page":1,"pages":2390}`
	client := newTestClientWithMockServer(data)
	quotes, _, err := client.Quotes()

	assert.Nil(t, err)
	assert.Equal(t, len(quotes), 1)
	assert.Equal(t, quotes[0].ID, "5cd96e05de30eff6ebcce7e9")
	assert.Equal(t, quotes[0].Character, "5cd99d4bde30eff6ebccfe9e")
}

func TestUnmarshalChapter(t *testing.T) {
	data := `{"docs":[{"_id":"6091b6d6d58360f988133b8b","chapterName":"A Long-expected Party","book":"5cf5805fb53e011a64671582"}],"total":62,"limit":1,"offset":0,"page":1,"pages":62}`

	client := newTestClientWithMockServer(data)
	chapters, _, err := client.Chapters()

	assert.Nil(t, err)
	assert.Equal(t, len(chapters), 1)
	assert.Equal(t, chapters[0].ID, "6091b6d6d58360f988133b8b")
	assert.Equal(t, chapters[0].Book, "5cf5805fb53e011a64671582")
}
