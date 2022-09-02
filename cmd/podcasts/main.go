package main

import (
	"net/http"
	"os"

	"github.com/CallumKerrEdwards/loggerrific"

	"github.com/CallumKerrEdwards/library-podcasts/internal/adapters/logrus"
	booksHTTPClient "github.com/CallumKerrEdwards/library-podcasts/internal/books/httpclient"
	mediaHTTPClient "github.com/CallumKerrEdwards/library-podcasts/internal/media/httpclient"
	"github.com/CallumKerrEdwards/library-podcasts/internal/podcasts/config"
	podcastsService "github.com/CallumKerrEdwards/library-podcasts/internal/podcasts/service"
	transportHttp "github.com/CallumKerrEdwards/library-podcasts/internal/transport/http"
)

// Run - sets the application.
func Run(logger loggerrific.Logger) error {
	logger.Infoln("Setting up Library Podcasts Server")

	cfg, err := config.New(logger)
	if err != nil {
		logger.WithError(err).Errorln("Config error")
		return err
	}

	booksClient := booksHTTPClient.NewBooksClient(cfg.Application.Dependencies.BooksAPIHost, http.DefaultClient, logger)
	mediaClient := mediaHTTPClient.NewMediaClient(cfg.Application.Dependencies.MediaAPIHost, http.DefaultClient, logger)

	mainPodcastService, err := podcastsService.New(cfg, booksClient, mediaClient, logger)
	if err != nil {
		logger.WithError(err).Errorln("service error")
		return err
	}

	httpHandler := transportHttp.NewHandler(cfg.Server, mainPodcastService, logger)
	if err = httpHandler.Serve(); err != nil {
		logger.WithError(err).Errorln("Server error")
		return err
	}

	return nil
}

func main() {
	logger := logrus.NewLogger()
	logger.SetLevelDebug()

	if err := Run(logger); err != nil {
		logger.WithError(err).Errorln("Error starting up Library Podcasts Server")
		os.Exit(1)
	}
}
