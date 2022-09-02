package responses

import "github.com/CallumKerrEdwards/library-podcasts/internal/adapters/dtos/books"

type AudiobooksDTO struct {
	Audiobooks []books.Book `json:"audiobooks"`
}
