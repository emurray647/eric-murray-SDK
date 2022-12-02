# Lord of the Rings SDK
 
This SDK is for accessing the [The One API to rule them all](http://the-one-api.dev). Note that this package requires Go 1.18.

## Table of Contents
- [Layout](#layout)
- [Usage](#usage)
    - [Client](#client)
    - [Filter](#filter)
- [Testing](#testing)
- [Future Improvements](#future-improvements)

## Layout

The layout of the project is as follows:
```
├── lotrsdk
│   ├── client.go
│   ├── client_test.go
│   ├── filter.go
│   ├── go.mod
│   ├── go.sum
│   └── model.go
└── README.md
```

A brief description of the files:
- `lotrsdk`: directory that contains all the source code
- `client.go`: defines the `Client` interface and implementation
- `client-test.go`: the unit tests for the `Client` interface
- `filter.go`: defines the `Filter` interface to enable filtering, pagination, and sorting
- `go.mod`: defines the module
- `go.sum`: generated fo file; do not edit
- `model.go`: defines the Go structs that correspond to the JSON responses
- `README.md`: description of the package

Note that the actual go module exists in the `lotrsdk` directory. This is so that
the repo is of the form `{name}-SDK`, but the package has a more resonable name
of `lotrsdk` (Lord of the Rings SDK).

## Usage

Before being able to use this sdk, you will need an access key to access the endpoints; 
you can get one [here](http://the-one-api.dev/sign-up).

To install this package, run
`go get github.com/emurray647/eric-murray-SDK/lotrsdk@v1.0.0`

In your `.go` files, import it as `"github.com/emurray647/eric-murray-SDK/lotrsdk"`

### Client

The `lotrsdk` package provides a `Client` interface one can access by `lotr.NewClient('<access-token'>)`.
The `Client` interface provides several methods that can be used to call the-one-api. Each method can take 0 or
more Filter objects ([see filter section](#filter)) to filter which records are selected, and returns 3 values: a slice of the searched for
values, a Status object, and an error type.  The error will be non-nil if there was any issue in the request, and 
the Status value just provides some information about the request (values like total items retrieved, limit, offset,
page, and total pages). The main return is the first value, which is a slice of the appropriate type (`Books` returns `[]Book`, 
`Movies` returns `[]Movie`, etc..). See `models.go` for a list of all the return types.  
All the methods of the `Client` interface are as follows:

| Method Name | Non-filter Parameters | Corresponding Endpoint | Description | Return Type |
| --- | --- | --- | --- | --- |
| `Books` | none | `/book` | List all of the books | `([]Book, Status, error)` |
| `ChapterFromBook` | `*Book` | `/book/{id}/chapter` | Request all chapters of the provided book | `([]Chapter, Status, error)` |
| `Movies` | none | `/movie` | Get all of the movies | `([]Movie, Status, error)` |
| `QuoteFromMovie` | `*Movie` | `/movie/{id}/quote` | Request all quotes from the provided movie |`([]Quote, Status, error)` |
| `Characters` | none | `/character` | Gets the list of characters |`([]Character, Status, error)` |
| `QuoteFromCharacter` | `*Character` | `/character/{id}/quote` | Request all quotes from the provided character |`([]Quote, Status, error)` |
| `Quotes` | none | `/quote` | Get the list of all quotes |`([]Quote, Status, error)` |
| `Chapters` | none | `/chapter` | Get the list of all chapters |`([]Chapter, Status, error)` |

So for example, to get the list of all movies, on could do

```
client := lotr.NewClient('<access-token>')
movies, status, err := client.Movies()
if err != nil {
    panic(err)
}
// do something with movies
```

### Filter

The `lotrsdk` package also provides several filtering options to use to select which records should be retrieved.
All of the `Client` methods can be passed zero or more `Filter` objects. The methods to create the `Filter` objects
are as follows:

- `BinaryFilter(key string, operator FilterCompareType, value string, values ...string)` \
Creates a `Filter` where key is compared to value using the provided operator. The operator can be one of FilterCompareEqual,
FilterCompareNotEqual, FilterCompareLessThan, FilterCompareGreaterThan, FilterCompareLessThanOrEqual, or FilterCompareGreaterThanOrEqual.
So for example, `BinaryFilter("name", FilterCompareEqual, "Gandalf")` will only select records whose `name` field is `"Gandalf"`, and
`BinaryFilter("budgetInMillions", FilterCompareLessThan, "100")` will only select records whose `budgetInMillions` field is strictly less
than `100`. \
\
There is the option to add more values when operator is either `FilterCompareEqual` or `FilterCompareNotEqual`, which is essentially checking the 
key against all the values. For instance, `BinaryFilter("name", FilterCompareEqual, "Gandalf", "Elrond")` will only select records where the `name`
field is either `"Gandalf"` or `"Elrond"`, while `BinaryFilter("name", FilterCompareNotEqual, "Gandalf", "Elrond")` will select all records
where the `name` field is not `"Gandalf"` or `"Elrond"`.

- `ExistFilter(key string)` \
Creates a `Filter` that only selects records where `key` exists as one of the fields.  For instance, `ExistFilter("wikiUrl")` will only select 
records that have the field `wikiUrl`.

- `NotExistFilter(key string)` \
Creates a `Filter` that only selects records where `key` does not exist as one of the fields.  For instance, `NotExistFilter("wikiUrl")` will
only select records that do not have a field named `wikiUrl`.

- `Sort(field string, order SortOrder)` \
Creates a `Filter` that instead of filtering, sorts the output based on the key `field`.  `order` can be either `SortOrderAscending` to sort
in ascending order, or `SortOrderDescending` to sort in descending order.

- `Limit(value int)` \
Creates a `Filter` that instead of filtering, limits the output to `value` records.  This is equivalent to a `limit={value}` query parameter.

- `Page(value int)` \
Creates a `Filter` that instead of filtering, selects the `value`th page. This is equivalent to a `page={value}` query parameter.

- `Offset(value int)` \
Creates a `Filter` that instead of filtering, skips the first `value`th records before selecting. This is equivalent
to a `offset={value}` query parameter.

Additionally, there is a convenience function `MergeFilters(filters ...Filter)` that returns a `Filter` which combines all the input `Filter`s.

As an example on how to use filters, let us say we want to find 5 quotes by a character named Gandalf:

```
client := NewClient("<access-token>")

// create a filter for all records with name="Gandalf"
filter := BinaryFilter("name", FilterCompareEqual, "Gandalf")

characters, _, err := client.Characters(filter) // ignore the status
if err != nil {
    panic(err)
} else if len(characters) < 0 {
    panic("Could not find character Gandalf")
}

// now use the character we found to look for quotes from Gandalf
gandalfCharacter := characters[0]
quotes, _, err := client.QuoteFromCharacter(&gandalfCharacter, Limit(5)) // only grab 5
if err != nil {
    panic(err)
}

for _, quote := range quotes {
    fmt.Println(quote.Dialog)
}
/*
Prints out:

Now come the days of the King. May they be blessed.
Hobbits!
Be careful. Even in defeat, Saruman is dangerous.
No, we need him alive. We need him to talk.
Your treachery has already cost many lives. Thousands more are now at risk. But you could save them Saruman. You were deep in the enemy's counsel.
*/
```

For a more complicated example, say we wanted to pick the first 5 alphetically named characters with a race of Elf or Maiar and non-brown hair 
who are also popular enough to have a wikiUrl:

```
client := NewClient("ULVbigXfcrP4otc04wVo")

// start by selecting race as Elf or Maiar
filter := BinaryFilter("race", FilterCompareEqual, "Elf", "Maiar")

// Now filter out records with hair that is not brown (or empty)
filter = MergeFilters(filter, BinaryFilter("hair", FilterCompareNotEqual, "brown", ""))

// no filter for records that have a wikiUrl field
filter = MergeFilters(filter, ExistFilter("wikiUrl"))

// sort by name
filter = MergeFilters(filter, Sort("name", SortOrderAscending))

// only get 5
filter = MergeFilters(filter, Limit(5))

characters, _, err := client.Characters(filter) // ignore the status
if err != nil {
    panic(err)
}

for _, character := range characters {
    fmt.Printf("name: %s, hair: %s, wiki: %s\n", character.Name, character.Hair, character.WikiURL)
}
/*
Prints out:

name: Aegnor, hair: Golden, wiki: http://lotr.wikia.com//wiki/Aegnor
name: Amras, hair: Dark red, wiki: http://lotr.wikia.com//wiki/Amras
name: Amrod, hair: Dark red, wiki: http://lotr.wikia.com//wiki/Amrod
name: Amroth, hair: Golden, wiki: http://lotr.wikia.com//wiki/Amroth
name: Angrod, hair: Golden, wiki: http://lotr.wikia.com//wiki/Angrod
*/
```

## Testing

Unit test can be run from the `lotrsdk/` directory with `go test ./...`

## Future Improvements
- Better testing
    - As it is, all tests are in `lotrsdk/client_test.go` and consist of either calling methods on `Client`, catching the request,
    and verifying it has the correct path and params, or explicitly mocking what the server returns and verifying we
    unmarshal correctly.  It would be nice to have some non-unit tests (that run separately from the unit tests) that call
    the actual API to verify the whole system is running correctly.
    - This module also needs more tests to check error conditions: in cases we get bad data or if we can't connect to the server, 
    we should make sure that an error is returned rather than incorrect data.
- Better type system for filtering. Right now the `Filter` methods work mostly on strings; for instance, 
`BinaryFilter("age", FilterCompareGreaterThan, "50")`. It would be nice if we could use an `int` there instead of a `string`.
- More time and care should be spent with how this module deals with query parameters.  Go's `net/url` package has a type `Values`
that usually works well for dealing with query parameters.  However, it heavily favors all keys and values in the format `key=value`
with only an equal sign.  This doesn't work for our case, as we want to be able to handle inequalities as well.  The solution 
in this `lotrsdk` package was hacked together in a hurry; it could use some more time to be cleaned up and fleshed out.
- Better handling of regex.  This package currently only does regex if you explicitly pass in a regex string 
(eg: `BinaryFilter("name", FilterCompareEqual, "/foot/i")`). It would be nice to integrate this with Go's `regexp` package.
This is probably a major lift that would take some time to do.
- Better naming for `Filter`: the `Filter` interface morphed into all of filter, pagination, and sorting. A better name
for the interface would make it more clear that it does not just filter.