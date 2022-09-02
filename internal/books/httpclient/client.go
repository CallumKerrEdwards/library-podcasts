package httpclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/CallumKerrEdwards/library-podcasts/pkg/books"
	"github.com/CallumKerrEdwards/loggerrific"
)

var (
	errGettingAudiobooks = errors.New("cannot get audiobooks")
)

type Client struct {
	BooksAPIHost string
	*http.Client
	Log loggerrific.Logger
}

func NewBooksClient(booksAPIHost string, httpClient *http.Client, logger loggerrific.Logger) *Client {
	return &Client{
		BooksAPIHost: booksAPIHost,
		Client:       httpClient,
		Log:          logger,
	}
}

func (c *Client) GetAllAudiobooks(ctx context.Context) ([]books.Book, error) {

	responseBodyMap := make(map[string][]books.Book)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BooksAPIHost+"/audiobooks", nil)
	if err != nil {
		c.Log.WithError(err).Errorln("Error creating audiobooks request")
		return []books.Book{}, nil
	}
	request.Header.Set("Accept", "application/json")

	response, err := c.Client.Do(request)
	if err != nil {
		c.Log.WithError(err).Errorln("Error getting audiobooks")
		return []books.Book{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		c.Log.WithField("status_code", response.StatusCode).Errorln("Error getting audiobooks")
		return []books.Book{}, fmt.Errorf("%w: status code %d", errGettingAudiobooks, response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(&responseBodyMap); err != nil {
		c.Log.WithError(err).Errorln("Cannot decode response into audiobooks")
		return []books.Book{}, nil
	}

	return responseBodyMap["audiobooks"], nil
}
