package httpclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/CallumKerrEdwards/loggerrific"

	"github.com/CallumKerrEdwards/library-podcasts/internal/adapters/dtos/books"
	"github.com/CallumKerrEdwards/library-podcasts/internal/adapters/dtos/responses"
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

	var responseBody responses.AudiobooksDTO

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

	// buf := new(strings.Builder)
	// _, err = io.Copy(buf, response.Body)
	// // check errors
	// c.Log.Errorln(buf.String())

	if response.StatusCode != 200 {
		c.Log.WithField("status_code", response.StatusCode).Errorln("Error getting audiobooks")
		return []books.Book{}, fmt.Errorf("%w: status code %d", errGettingAudiobooks, response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		c.Log.Errorln(err)
		c.Log.WithError(err).Errorln("Cannot decode response into audiobooks")
		return []books.Book{}, err
	}

	return responseBody.Audiobooks, nil
}
