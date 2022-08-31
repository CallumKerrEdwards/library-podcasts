package service

import (
	"github.com/CallumKerrEdwards/loggerrific"
	"github.com/jbub/podcasts"
)

// type BooksClient interface {
// }

var (
	defaultRootPath = "/data/feeds"
)

type EpisodesClient interface {
	GetAllEpisodes() ([]*podcasts.Item, error)
}

// Service - provides all functions for accessing and modifying Books.
type Service struct {
	EpisodesClient EpisodesClient
	Log            loggerrific.Logger
	rootPath       string
}

func New(episodesClient EpisodesClient, logger loggerrific.Logger) (*Service, error) {
	svc := &Service{
		EpisodesClient: episodesClient,
		Log:            logger,
		rootPath:       defaultRootPath,
	}
	err := svc.initFeeds()
	if err != nil {
		return &Service{}, err
	}
	return svc, nil
}

func (s *Service) GetFeedsRoot() string {
	return s.rootPath
}
