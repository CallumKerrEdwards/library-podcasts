package books

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Book - representation of a book.
type Book struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Authors     []Person     `json:"authors"`
	Description *Description `json:"description,omitempty"`
	ReleaseDate *ReleaseDate `json:"releaseDate,omitempty"`
	Genres      []Genre      `json:"genres,omitempty"`
	Series      Series       `json:"series"`
	Audiobook   *Audiobook   `json:"audiobook,omitempty"`
}

// Person - represetation of a person, for example an author or audiobook narrator.
type Person struct {
	Forenames string `json:"forenames"`
	SortName  string `json:"sortName"`
}

// Series - representation of a series of books.
type Series struct {
	Sequence int    `json:"sequence"`
	Title    string `json:"title"`
}

func NewBook(title string, description *Description, authors []Person, releaseDate *ReleaseDate,
	genreList []Genre, series Series, audiobook *Audiobook) Book {
	return Book{
		ID:          uuid.New().String(),
		Title:       title,
		Authors:     authors,
		Description: description,
		ReleaseDate: releaseDate,
		Genres:      genreList,
		Series:      series,
		Audiobook:   audiobook,
	}
}

// Description - representation of a blurb of a book.
type Description struct {
	Text   string `json:"text"`
	Format Format `json:"format,omitempty"`
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

type ReleaseDate struct {
	time.Time
}

const releaseDateLayout = "2006-01-02"

func NewReleaseDate(date string) (*ReleaseDate, error) {
	layout := "2006-01-02"

	datetime, err := time.Parse(layout, date)
	if err != nil {
		return nil, err
	}

	return &ReleaseDate{Time: datetime}, nil
}

func (d *ReleaseDate) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		d.Time = time.Time{}
		return
	}

	d.Time, err = time.Parse(releaseDateLayout, s)

	return
}

func (d *ReleaseDate) MarshalJSON() ([]byte, error) {
	if d.Time.UnixNano() == nilTime {
		return []byte("null"), nil
	}

	return []byte(fmt.Sprintf("%q", d.Time.Format(releaseDateLayout))), nil
}

var nilTime = (time.Time{}).UnixNano()

func (d *ReleaseDate) IsSet() bool {
	return d.UnixNano() != nilTime
}

type Audiobook struct {
	AudiobookMediaID             string   `json:"audiobookMediaId"`
	Narrators                    []Person `json:"narrators"`
	CoverImageMediaID            string   `json:"coverImageMediaId"`
	SupplimentaryMaterialMediaID string   `json:"supplimentaryMaterialMediaId,omitempty"`
}

func (a Audiobook) GetNarrator() string {
	return GetPersonsString(a.Narrators)
}
