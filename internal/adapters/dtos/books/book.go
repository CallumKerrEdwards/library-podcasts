package books

import (
	"fmt"
	"strings"
)

// Book - representation of a book.
type Book struct {
	ID          string       `json:"id" pact:"example=fc763eba-0905-41c5-a27f-3934ab26786c,regex=^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"`
	Title       string       `json:"title" pact:"example=The Final Empire"`
	Authors     []Person     `json:"authors"`
	Description *Description `json:"description,omitempty"`
	ReleaseDate string       `json:"releaseDate,omitempty" pact:"example=2000-01-01,regex=^\\d{4}-\\d{2}-\\d{2}$"`
	Genres      []string     `json:"genres,omitempty"`
	Series      Series       `json:"series"`
	Audiobook   *Audiobook   `json:"audiobook,omitempty"`
}

// Person - represetation of a person, for example an author or audiobook narrator.
type Person struct {
	Forenames string `json:"forenames" pact:"example=Brandon"`
	SortName  string `json:"sortName" pact:"example=Sanderson"`
}

// Series - representation of a series of books.
type Series struct {
	Sequence int    `json:"sequence" pact:"example=1"`
	Title    string `json:"title" pact:"example=Mistborn"`
}

// Description - representation of a blurb of a book.
type Description struct {
	Text   string `json:"text" pact:"example=<h1>The Final Empire</h1><p>Ash fell from the sky.</p>"`
	Format string `json:"format,omitempty" pact:"example=HTML,regex=^(HTML|Markdown|Plain)$"`
}

func (b *Book) GetAuthor() string {
	return GetPersonsString(b.Authors)
}

func GetPersonsString(p []Person) string {
	switch len(p) {
	case 0:
		return ""
	case 1:
		return p[0].GetPersonString()
	default:
		var personStrs []string
		for _, person := range p {
			personStrs = append(personStrs, person.GetPersonString())
		}

		return fmt.Sprintf("%s & %s", strings.Join(personStrs[:len(personStrs)-1], ", "), personStrs[len(personStrs)-1])
	}
}

func (p Person) GetPersonString() string {
	return fmt.Sprintf("%s %s", p.Forenames, p.SortName)
}

type Audiobook struct {
	AudiobookMediaID             string   `json:"audiobookMediaId" pact:"example=44864ec9-77cb-4c76-8d55-c22321a8b51c,regex=^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"`
	Narrators                    []Person `json:"narrators"`
	CoverImageMediaID            string   `json:"coverImageMediaId" pact:"example=2d5ba852-87f6-4902-8a3f-49cfa2913352,regex=^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"`
	SupplimentaryMaterialMediaID string   `json:"supplimentaryMaterialMediaId,omitempty" pact:"example=fc763eba-0905-41c5-a27f-3934ab26786c,regex=^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"`
}

func (a Audiobook) GetNarrator() string {
	return GetPersonsString(a.Narrators)
}
