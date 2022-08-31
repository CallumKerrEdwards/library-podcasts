package main

import (
	"os"

	"github.com/CallumKerrEdwards/loggerrific"

	"github.com/CallumKerrEdwards/library-podcasts/internal/adapters/logrus"

	"github.com/CallumKerrEdwards/library-podcasts/internal/episodes/dummyclient"
	podcastsService "github.com/CallumKerrEdwards/library-podcasts/internal/podcasts/service"
	transportHttp "github.com/CallumKerrEdwards/library-podcasts/internal/transport/http"
)

// Run - sets the application.
func Run(logger loggerrific.Logger) error {
	logger.Infoln("Setting up Library Podcasts Server")

	episodesClient := &dummyclient.DummyEpisodesClient{}

	mainPodcastService, err := podcastsService.New(episodesClient, logger)
	if err != nil {
		logger.WithError(err).Errorln("service error")
		return err
	}

	httpHandler := transportHttp.NewHandler(mainPodcastService, logger)
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
