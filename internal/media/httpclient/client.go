package httpclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/CallumKerrEdwards/loggerrific"

	"github.com/CallumKerrEdwards/library-podcasts/pkg/media"
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

func (c *Client) getMediaPath() string {
	return fmt.Sprintf("%s/cms/v1/media", c.MediaAPIHost)
}

func (c *Client) GetMedia(ctx context.Context, id string) (media.Media, error) {
	var responseBody media.Media

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", c.getMediaPath(), id), nil)
	if err != nil {
		c.Log.WithError(err).Errorln("Error creating media request")
		return media.Media{}, err
	}
	request.Header.Set("Accept", "application/json")

	response, err := c.Client.Do(request)
	if err != nil {
		c.Log.WithError(err).Errorln("Error getting media")
		return media.Media{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		c.Log.WithField("status_code", response.StatusCode).Errorln("Error getting media")
		return media.Media{}, fmt.Errorf("%w: status code %d", errGettingMedia, response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		c.Log.WithError(err).Errorln("Cannot decode response into media")
		return media.Media{}, err
	}

	return responseBody, nil
}

type pathResponse struct {
	Path string `json:"path" pact:"example=/media/audio/9136ee1d-f0ba-428e-9092-ad64f3ab98863-Pride%20and%20Prejudice.m4b"`
}

func (c *Client) GetPath(ctx context.Context, id string) (string, error) {

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s/path", c.getMediaPath(), id), nil)
	if err != nil {
		c.Log.WithError(err).Errorln("Error creating path request")
		return "", err
	}
	request.Header.Set("Accept", "application/json")

	response, err := c.Client.Do(request)
	if err != nil {
		c.Log.WithError(err).Errorln("Error getting path")
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		c.Log.WithField("status_code", response.StatusCode).Errorln("Error getting path")
		return "", fmt.Errorf("%w: status code %d", errGettingMedia, response.StatusCode)
	}

	var responseBody pathResponse
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		c.Log.WithError(err).Errorln("Cannot decode response")
		return "", err
	}

	return responseBody.Path, nil
}
