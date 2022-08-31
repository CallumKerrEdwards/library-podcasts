package service

import (
	"os"
	"path/filepath"

	"github.com/jbub/podcasts"
)

func (s *Service) initFeeds() error {

	err := os.MkdirAll(s.rootPath, os.ModePerm)
	if err != nil {
		s.Log.WithError(err).Errorln("Cannot create path")
		return err
	}

	// initialize the podcast
	p := &podcasts.Podcast{
		Title:       "My podcast",
		Description: "This is my very simple podcast.",
		Language:    "EN",
		Link:        "http://www.example-podcast.com/my-podcast",
		Copyright:   "2015 My podcast copyright",
	}

	allItems, err := s.EpisodesClient.GetAllEpisodes()
	if err != nil {
		return err
	}

	for _, item := range allItems {
		p.AddItem(item)
	}

	// get podcast feed, you can pass options to customize it
	feed, err := p.Feed(
		podcasts.Author("Author Name"),
		podcasts.Block,
		podcasts.Explicit,
		podcasts.Complete,
		podcasts.NewFeedURL("http://www.example-podcast.com/new-feed-url"),
		podcasts.Subtitle("This is my very simple podcast subtitle."),
		podcasts.Summary("This is my very simple podcast summary."),
		podcasts.Owner("Podcast Owner", "owner@example-podcast.com"),
		podcasts.Image("http://www.example-podcast.com/my-podcast.jpg"),
	)

	// handle error
	if err != nil {
		s.Log.WithError(err).Errorln("Cannot create podcast feed")
		return err
	}

	mainFeedFilepath := filepath.Join(s.rootPath, "main.rss")
	file, err := os.Create(mainFeedFilepath)
	defer file.Close()

	if err != nil {
		s.Log.WithError(err).Errorln("Cannot create main.rss")
		return err
	}

	err = feed.Write(file)
	if err != nil {
		s.Log.WithError(err).Errorln("Cannot write podcast feed to main.rss")
		return err
	}

	s.Log.Infoln("Initialised podcast feed at", mainFeedFilepath)

	return nil
}
