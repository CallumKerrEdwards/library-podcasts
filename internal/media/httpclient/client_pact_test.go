//go:build contract
// +build contract

package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/CallumKerrEdwards/loggerrific/tlogger"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/pact-foundation/pact-go/v2/models"
	"github.com/stretchr/testify/assert"

	"github.com/CallumKerrEdwards/library-podcasts/pkg/media"
)

func TestMediaPathContract(t *testing.T) {

	existingId := "9136ee1d-f0ba-428e-9092-ad64f3ab98863"
	notFoundId := "2"

	mockProvider, err := consumer.NewV3Pact(consumer.MockHTTPProviderConfig{
		Consumer: "LibraryPodcastsMediaEndpointConsumer",
		Provider: "LibraryCMSMediaEndpointProvider",
	})
	assert.NoError(t, err)

	err = mockProvider.
		AddInteraction().
		GivenWithParameter(models.ProviderState{
			Name: "Media with id exists",
			Parameters: map[string]interface{}{
				"id":    existingId,
				"title": "Pride and Prejudice",
			},
		}).
		UponReceiving("A request for a media item with specific id").
		WithRequestPathMatcher("GET", matchers.S("/cms/v1/media/"+existingId)).
		WillRespondWith(200, func(b *consumer.V3ResponseBuilder) {
			b.
				Header("Content-Type", matchers.S("application/json")).
				BodyMatch(&media.Media{})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			testClient := NewMediaClient(fmt.Sprintf("http://%s:%d", config.Host, config.Port), http.DefaultClient, tlogger.NewTLogger(t))

			med, err := testClient.GetMedia(context.Background(), existingId)

			assert.NoError(t, err)
			assert.Equal(t, "Pride and Prejudice", med.Title)
			assert.Equal(t, int64(3869609), med.Size)

			return err
		})
	assert.NoError(t, err)

	err = mockProvider.
		AddInteraction().
		GivenWithParameter(models.ProviderState{
			Name: "Media with id does not exist",
			Parameters: map[string]interface{}{
				"id": notFoundId,
			},
		}).
		UponReceiving("A request for a media item with non-existent id").
		WithRequestPathMatcher("GET", matchers.S("/cms/v1/media/"+notFoundId)).
		WillRespondWith(404, func(b *consumer.V3ResponseBuilder) {
			b.
				Header("Content-Type", matchers.S("application/json"))
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			testClient := NewMediaClient(fmt.Sprintf("http://%s:%d", config.Host, config.Port), http.DefaultClient, tlogger.NewTLogger(t))

			_, err := testClient.GetMedia(context.Background(), notFoundId)
			if assert.Error(t, err) {
				assert.EqualError(t, err, "cannot get media: status code 404")
				return nil
			}

			t.Fail()
			return nil
		})
	assert.NoError(t, err)

	err = mockProvider.
		AddInteraction().
		GivenWithParameter(models.ProviderState{
			Name: "Media with id exists",
			Parameters: map[string]interface{}{
				"id":    existingId,
				"title": "Pride and Prejudice",
			},
		}).
		UponReceiving("A request for path of media item with specific id").
		WithRequestPathMatcher("GET", matchers.S("/cms/v1/media/"+existingId+"/path")).
		WillRespondWith(200, func(b *consumer.V3ResponseBuilder) {
			b.
				Header("Content-Type", matchers.S("application/json")).
				BodyMatch(&pathResponse{})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			testClient := NewMediaClient(fmt.Sprintf("http://%s:%d", config.Host, config.Port), http.DefaultClient, tlogger.NewTLogger(t))

			path, err := testClient.GetPath(context.Background(), existingId)

			assert.NoError(t, err)
			assert.Equal(t, "/media/audio/"+existingId+"-Pride%20and%20Prejudice.m4b", path)

			return err
		})
	assert.NoError(t, err)

	err = mockProvider.
		AddInteraction().
		GivenWithParameter(models.ProviderState{
			Name: "Media with id does not exist",
			Parameters: map[string]interface{}{
				"id": notFoundId,
			},
		}).
		UponReceiving("A request for path of media item with non-existent id").
		WithRequestPathMatcher("GET", matchers.S("/cms/v1/media/"+notFoundId+"/path")).
		WillRespondWith(404, func(b *consumer.V3ResponseBuilder) {
			b.
				Header("Content-Type", matchers.S("application/json"))
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			testClient := NewMediaClient(fmt.Sprintf("http://%s:%d", config.Host, config.Port), http.DefaultClient, tlogger.NewTLogger(t))

			_, err := testClient.GetPath(context.Background(), notFoundId)
			if assert.Error(t, err) {
				assert.EqualError(t, err, "cannot get media: status code 404")
				return nil
			}

			t.Fail()
			return nil
		})
	assert.NoError(t, err)
}
