package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CallumKerrEdwards/loggerrific/tlogger"
	"github.com/stretchr/testify/assert"
)

const (
	audiobooksResponse = `
	{
		"audiobooks": [
		  {
			"id": "b05b79e2-17ff-4896-848d-a5a6a5e78066",
			"title": "The Final Empire",
			"authors": [
			  {
				"forenames": "Brandon",
				"sortName": "Sanderson"
			  }
			],
			"series": {
			  "sequence": 1,
			  "title": "Mistborn"
			},
			"audiobook": {
			  "audiobookMediaId": "0001",
			  "narrators": [
				{
				  "forenames": "Michael",
				  "sortName": "Kramer"
				}
			  ]
			}
		  },
		  {
			"id": "53057ac4-9aa1-4369-a840-366e11a0d156",
			"title": "The Well of Ascension",
			"authors": [
			  {
				"forenames": "Brandon",
				"sortName": "Sanderson"
			  }
			],
			"series": {
			  "sequence": 2,
			  "title": "Mistborn"
			},
			"audiobook": {
			  "audiobookMediaId": "0002",
			  "narrators": [
				{
				  "forenames": "Michael",
				  "sortName": "Kramer"
				}
			  ]
			}
		  },
		  {
			"id": "50fd142a-8c35-4460-892e-efc0f8eb81d4",
			"title": "The Hero of Ages",
			"authors": [
			  {
				"forenames": "Brandon",
				"sortName": "Sanderson"
			  }
			],
			"series": {
			  "sequence": 3,
			  "title": "Mistborn"
			},
			"audiobook": {
			  "audiobookMediaId": "0003",
			  "narrators": [
				{
				  "forenames": "Michael",
				  "sortName": "Kramer"
				}
			  ]
			}
		  }
		]
	  }
`
)

func TestBooksClient(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/audiobooks")
		rw.Write([]byte(audiobooksResponse))
		rw.Header().Add("Content-Type", "application/json")
	}))
	defer server.Close()

	testClient := NewBooksClient(server.URL, server.Client(), tlogger.NewTLogger(t))

	allAudiobooks, err := testClient.GetAllAudiobooks(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, "b05b79e2-17ff-4896-848d-a5a6a5e78066", allAudiobooks[0].ID)
	assert.Equal(t, "53057ac4-9aa1-4369-a840-366e11a0d156", allAudiobooks[1].ID)
	assert.Equal(t, "50fd142a-8c35-4460-892e-efc0f8eb81d4", allAudiobooks[2].ID)
}
