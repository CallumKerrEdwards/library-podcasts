package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/CallumKerrEdwards/library-podcasts/internal/adapters/dtos/responses"
	"github.com/CallumKerrEdwards/loggerrific/tlogger"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/log"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

type S = matchers.S

var Decimal = matchers.Decimal
var Integer = matchers.Integer
var Equality = matchers.Equality
var Includes = matchers.Includes
var FromProviderState = matchers.FromProviderState
var EachKeyLike = matchers.EachKeyLike
var ArrayContaining = matchers.ArrayContaining
var ArrayMinMaxLike = matchers.ArrayMinMaxLike
var ArrayMaxLike = matchers.ArrayMaxLike
var DateGenerated = matchers.DateGenerated
var TimeGenerated = matchers.TimeGenerated
var DateTimeGenerated = matchers.DateTimeGenerated

var Like = matchers.Like
var EachLike = matchers.EachLike
var Term = matchers.Term
var Regex = matchers.Regex
var HexValue = matchers.HexValue
var Identifier = matchers.Identifier
var IPAddress = matchers.IPAddress
var IPv6Address = matchers.IPv6Address
var Timestamp = matchers.Timestamp
var Date = matchers.Date
var Time = matchers.Time
var UUID = matchers.UUID
var ArrayMinLike = matchers.ArrayMinLike

type Map = matchers.MapMatcher

// type GetAudiobooksResponse struct {
// 	Audiobooks []books.Book `json:"audiobooks", `
// }

func TestUserAPIClient(t *testing.T) {
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
		WithRequest("Get", "/audiobooks").
		WillRespondWith(200, func(b *consumer.V3ResponseBuilder) {
			b.
				Header("Content-Type", S("application/json")).
				BodyMatch(responses.AudiobooksDTO{})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			// Act: test our API client behaves correctly
			// Initialise the API client and point it at the Pact mock server

			testClient := NewBooksClient(fmt.Sprintf("http://%s:%d", config.Host, config.Port), http.DefaultClient, tlogger.NewTLogger(t))

			// Execute the API client
			audiobooks, err := testClient.GetAllAudiobooks(context.Background())

			// Assert: check the result
			assert.NoError(t, err)
			assert.NotEmpty(t, audiobooks)
			assert.Equal(t, "fc763eba-0905-41c5-a27f-3934ab26786c", audiobooks[0].ID)

			return err
		})
	assert.NoError(t, err)
}
