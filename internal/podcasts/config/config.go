package config

import (
	"path/filepath"
	"strings"

	"github.com/CallumKerrEdwards/loggerrific"
	"github.com/spf13/viper"
)

type Config struct {
	Application Application
	Podcast     Podcast
	Server      Server
}

type Application struct {
	Host         string
	Dependencies Dependencies
}

type Dependencies struct {
	BooksAPIHost string
	MediaAPIHost string
}

type Server struct {
	PodcastFeedsPathPrefix string
	PodcastFeedFilename    string
}

type Podcast struct {
	Title           string
	Subtitle        string
	Description     string
	Explicit        bool
	Copyright       string
	Language        string
	Author          string
	BlockFromITunes bool
}

func New(logger loggerrific.Logger) (Config, error) {

	viper.SetDefault("Podcast.Title", "Audiobooks")
	viper.SetDefault("Podcast.Explicit", true)
	viper.SetDefault("Podcast.BlockFromITunes", true)
	viper.SetDefault("Server.PodcastFeedsPathPrefix", "/feed")
	viper.SetDefault("Server.PodcastFeedFilename", "feed.rss")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("libpodcasts")
	viper.BindEnv("config_path")

	viper.BindEnv("Application.Host")
	viper.BindEnv("Application.Dependencies.BooksAPIHost")
	viper.BindEnv("Application.Dependencies.MediaAPIHost")

	viper.BindEnv("Podcast.Title")
	viper.BindEnv("Podcast.Subtitle")
	viper.BindEnv("Podcast.Description")
	viper.BindEnv("Podcast.Explicit")
	viper.BindEnv("Podcast.Copyright")
	viper.BindEnv("Podcast.Language")
	viper.BindEnv("Podcast.Author")
	viper.BindEnv("Podcast.BlockFromITunes")

	viper.AutomaticEnv()

	pathToConfig := viper.GetString("config_path")

	if !filepath.IsAbs(pathToConfig) {
		logger.Infoln("No valid config path found from environment variable LIBPODCASTS_CONFIG_PATH, reading config from environment variables only")
	} else {
		viper.SetConfigFile(pathToConfig)
		err := viper.ReadInConfig()
		if err != nil {
			logger.WithError(err).Errorln("Cannot read config from file")
			return Config{}, err
		}
	}

	var podcastConfig Config
	err := viper.Unmarshal(&podcastConfig)

	return podcastConfig, err
}
