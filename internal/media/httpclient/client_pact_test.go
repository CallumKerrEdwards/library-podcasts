package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/CallumKerrEdwards/loggerrific/tlogger"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/log"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/pact-foundation/pact-go/v2/models"
	"github.com/stretchr/testify/assert"
)

func TestUserAPIClient(t *testing.T) {

	id := "9136ee1d-f0ba-428e-9092-ad64f3ab98863"

	log.SetLogLevel("TRACE")
	mockProvider, err := consumer.NewV3Pact(consumer.MockHTTPProviderConfig{
		Consumer: "PactGoV3Consumer",
		Provider: "V3Provider",
	})
	assert.NoError(t, err)

	// Set up our expected interactions.
	err = mockProvider.
		AddInteraction().
		Given("Some books with audiobook artefacts exist").
		UponReceiving("A request for all audiobooks").
		GivenWithParameter(models.ProviderState{
			Name: "Media with id exists",
			Parameters: map[string]interface{}{
				"id": id,
			},
		}).
		WithRequestPathMatcher("GET", matchers.S("/"+id+"/path")).
		WillRespondWith(200, func(b *consumer.V3ResponseBuilder) {
			b.
				Header("Content-Type", matchers.S("application/json")).
				BodyMatch(b.
					JSONBody(matchers.Map{
						"path": matchers.S("/media/audio/" + id + "-Pride%20and%20Prejudice.m4b"),
					},
					))
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			testClient := NewMediaClient(fmt.Sprintf("http://%s:%d", config.Host, config.Port), http.DefaultClient, tlogger.NewTLogger(t))

			// Execute the API client
			path, err := testClient.GetPath(context.Background(), id)

			// Assert: check the result
			assert.NoError(t, err)
			assert.Equal(t, "/media/audio/"+id+"-Pride%20and%20Prejudice.m4b", path)

			return err
		})
	assert.NoError(t, err)
}
