package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/CallumKerrEdwards/podcasts"
)

func (s *Service) InitFeeds(ctx context.Context) error {

	err := os.MkdirAll(s.rootPath, os.ModePerm)
	if err != nil {
		s.Log.WithError(err).Errorln("Cannot create path")
		return err
	}

	// initialize the podcast
	p := &podcasts.Podcast{
		Title:       s.Config.Podcast.Title,
		Description: s.Config.Podcast.Description,
		Language:    s.Config.Podcast.Language,
		Link:        s.Config.Application.Host,
		Copyright:   s.Config.Podcast.Copyright,
	}

	allAudiobooks, err := s.AudiobookClient.GetAllAudiobooks(ctx)
	if err != nil {
		return err
	}

	for _, audiobook := range allAudiobooks {

		audiobookMedia, err := s.MediaClient.GetMedia(ctx, audiobook.Audiobook.AudiobookMediaID)
		if err != nil {
			s.Log.WithError(err).Errorln("Cannot get audiobook media for book with id", audiobook.ID)
			continue
		}

		audiobookPath, err := s.MediaClient.GetPath(ctx, audiobook.Audiobook.AudiobookMediaID)
		if err != nil {
			s.Log.WithError(err).Errorln("Cannot get audiobook path for book with id", audiobook.ID)
			continue
		}

		fullAudiobookFileLocation := s.Config.Application.Host + audiobookPath

		podcastItem := &podcasts.Item{
			Title:          audiobook.ID,
			GUID:           fullAudiobookFileLocation,
			PubDate:        podcasts.NewPubDate(audiobook.ReleaseDate.Time),
			Description:    &podcasts.CDATAText{Value: fmt.Sprintf("%s by %s", audiobook.Title, audiobook.GetAuthor())},
			ContentEncoded: &podcasts.CDATAText{Value: ""},
			Enclosure: &podcasts.Enclosure{
				URL:    fullAudiobookFileLocation,
				Length: fmt.Sprint(audiobookMedia.Size),
				Type:   audiobookMedia.MIMEType,
			},
		}

		p.AddItem(podcastItem)
	}

	// get podcast feed, you can pass options to customize it
	feed, err := p.Feed(
		podcasts.Author(s.Config.Podcast.Author),
		podcasts.Block,
		podcasts.Explicit,
		podcasts.Subtitle(s.Config.Podcast.Subtitle),
		podcasts.Summary(s.Config.Podcast.Description),
		// podcasts.Image("http://www.example-podcast.com/my-podcast.jpg"),
	)

	// handle error
	if err != nil {
		s.Log.WithError(err).Errorln("Cannot create podcast feed")
		return err
	}

	mainFeedFilepath := filepath.Join(s.rootPath, s.Config.Server.PodcastFeedFilename)
	file, err := os.Create(mainFeedFilepath)
	defer file.Close()

	if err != nil {
		s.Log.WithError(err).Errorln("Cannot create", s.Config.Server.PodcastFeedFilename)
		return err
	}

	err = feed.Write(file)
	if err != nil {
		s.Log.WithError(err).Errorln("Cannot write podcast feed to", s.Config.Server.PodcastFeedFilename)
		return err
	}

	s.Log.Infoln("Initialised podcast feed at", mainFeedFilepath)

	return nil
}
