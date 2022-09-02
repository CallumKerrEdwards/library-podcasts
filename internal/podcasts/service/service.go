package service

import (
	"context"

	"github.com/CallumKerrEdwards/loggerrific"

	"github.com/CallumKerrEdwards/library-podcasts/internal/podcasts/config"
	"github.com/CallumKerrEdwards/library-podcasts/pkg/books"
	"github.com/CallumKerrEdwards/library-podcasts/pkg/media"
)

// type BooksClient interface {
// }

var (
	defaultRootPath = "/data/feeds"
)

type AudiobookClient interface {
	GetAllAudiobooks(context.Context) ([]books.Book, error)
}

type MediaClient interface {
	GetMedia(ctx context.Context, id string) (media.Media, error)
	GetPath(ctx context.Context, id string) (string, error)
}

// Service - provides all functions for accessing and modifying Books.
type Service struct {
	Config          config.Config
	AudiobookClient AudiobookClient
	MediaClient     MediaClient
	Log             loggerrific.Logger
	rootPath        string
}

func New(config config.Config, audiobooksClient AudiobookClient, mediaClient MediaClient, logger loggerrific.Logger) (*Service, error) {
	svc := &Service{
		Config:          config,
		AudiobookClient: audiobooksClient,
		MediaClient:     mediaClient,
		Log:             logger,
		rootPath:        defaultRootPath,
	}
	return svc, nil
}

func (s *Service) GetFeedsRoot() string {
	return s.rootPath
}
