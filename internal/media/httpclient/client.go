package httpclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/CallumKerrEdwards/loggerrific"

	"github.com/CallumKerrEdwards/library-podcasts/internal/adapters/dtos/media"
)

var (
	errGettingMedia = errors.New("cannot get media")
)

type Client struct {
	MediaAPIHost string
	*http.Client
	Log loggerrific.Logger
}

func NewMediaClient(MediaAPIHost string, httpClient *http.Client, logger loggerrific.Logger) *Client {
	return &Client{
		MediaAPIHost: MediaAPIHost,
		Client:       httpClient,
		Log:          logger,
	}
}

func (c *Client) GetMedia(ctx context.Context, id string) (media.Media, error) {
	var responseBody media.Media

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", c.MediaAPIHost, id), nil)
	if err != nil {
		c.Log.WithError(err).Errorln("Error creating media request")
		return media.Media{}, nil
	}
	request.Header.Set("Accept", "application/json")

	response, err := c.Client.Do(request)
	if err != nil {
		c.Log.WithError(err).Errorln("Error getting media")
		return media.Media{}, nil
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		c.Log.WithField("status_code", response.StatusCode).Errorln("Error getting media")
		return media.Media{}, fmt.Errorf("%w: status code %d", err, response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		c.Log.WithError(err).Errorln("Cannot decode response into media")
	}

	return responseBody, nil
}

// type pathResponse struct {
// 	ID string `json:"id"`
// 	Path strign
// }

func (c *Client) GetPath(ctx context.Context, id string) (string, error) {
	responseBodyMap := make(map[string]string)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s/path", c.MediaAPIHost, id), nil)
	if err != nil {
		c.Log.WithError(err).Errorln("Error creating path request")
		return "", nil
	}
	request.Header.Set("Accept", "application/json")

	response, err := c.Client.Do(request)
	if err != nil {
		c.Log.WithError(err).Errorln("Error getting path")
		return "", nil
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		c.Log.WithField("status_code", response.StatusCode).Errorln("Error getting path")
		return "", fmt.Errorf("%w: status code %d", err, response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(&responseBodyMap); err != nil {
		c.Log.Errorln()
		c.Log.WithError(err).Errorln("Cannot decode response")
		return "", err
	}

	return responseBodyMap["path"], nil
}
