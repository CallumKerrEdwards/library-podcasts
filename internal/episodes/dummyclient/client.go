package dummyclient

import (
	"sort"
	"time"

	"github.com/jbub/podcasts"
)

type DummyEpisodesClient struct {
}

func (c *DummyEpisodesClient) GetAllEpisodes() ([]*podcasts.Item, error) {
	var items []*podcasts.Item

	items = append(items, &podcasts.Item{
		Title:    "Episode 2",
		GUID:     "http://www.example-podcast.com/my-podcast/2/episode-two",
		PubDate:  podcasts.NewPubDate(time.Now()),
		Duration: podcasts.NewDuration(time.Second * 320),
		Enclosure: &podcasts.Enclosure{
			URL:    "http://www.example-podcast.com/my-podcast/2/episode.mp3",
			Length: "46732",
			Type:   "MP3",
		},
	})

	items = append(items, &podcasts.Item{
		Title:    "Episode 1",
		GUID:     "http://www.example-podcast.com/my-podcast/1/episode-one",
		PubDate:  podcasts.NewPubDate(time.Now().Local().Add(time.Hour * time.Duration(-14*24))),
		Duration: podcasts.NewDuration(time.Second * 230),
		Enclosure: &podcasts.Enclosure{
			URL:    "http://www.example-podcast.com/my-podcast/1/episode.mp3",
			Length: "12312",
			Type:   "MP3",
		},
	})

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].PubDate.Time.After(items[j].PubDate.Time)
	})

	return items, nil
}
