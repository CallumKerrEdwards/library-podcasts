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
	mediaID       = "c3bbe250-5532-454f-94cb-47ee0e0588d9"
	mediaResponse = `
	{
		"id": "c3bbe250-5532-454f-94cb-47ee0e0588d9",
		"title": "Pride and Prejudice",
		"mimeType": "audio/mp4a-latm",
		"extension": ".m4b",
		"size": 3869609
	}
`
	pathResponseBody = `
	{
		"id": "c3bbe250-5532-454f-94cb-47ee0e0588d9",
		"path": "/media/audio/c3bbe250-5532-454f-94cb-47ee0e0588d9-Pride%20and%20Prejudice.m4b"
	}
`
)

func TestBooksClient_GetMedia(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/cms/v1/media/"+mediaID, req.URL.String())
		rw.Write([]byte(mediaResponse))
		rw.Header().Add("Content-Type", "application/json")
	}))
	defer server.Close()

	testClient := NewMediaClient(server.URL, server.Client(), tlogger.NewTLogger(t))

	fetchedMedia, err := testClient.GetMedia(context.Background(), mediaID)
	assert.Nil(t, err)
	assert.Equal(t, mediaID, fetchedMedia.ID)
	assert.Equal(t, "audio/mp4a-latm", fetchedMedia.MIMEType)
	assert.Equal(t, int64(3869609), fetchedMedia.Size)
}

func TestBooksClient_GetPath(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/cms/v1/media/"+mediaID+"/path", req.URL.String())
		rw.Write([]byte(pathResponseBody))
		rw.Header().Add("Content-Type", "application/json")
	}))
	defer server.Close()

	testClient := NewMediaClient(server.URL, server.Client(), tlogger.NewTLogger(t))

	path, err := testClient.GetPath(context.Background(), mediaID)
	assert.NoError(t, err)
	assert.Equal(t, "/media/audio/c3bbe250-5532-454f-94cb-47ee0e0588d9-Pride%20and%20Prejudice.m4b", path)

}
