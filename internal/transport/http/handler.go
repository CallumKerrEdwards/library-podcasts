package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/CallumKerrEdwards/library-podcasts/internal/podcasts/config"
	"github.com/CallumKerrEdwards/loggerrific"
	"github.com/gorilla/mux"
)

var (
	address = "0.0.0.0:8080"
)

type FeedService interface {
	GetFeedsRoot() string
}

type Handler struct {
	Router  *mux.Router
	Service FeedService
	Config  config.Server
	Server  *http.Server
	Log     loggerrific.Logger
}

func NewHandler(cfg config.Server, service FeedService, logger loggerrific.Logger) *Handler {
	h := &Handler{
		Service: service,
		Config:  cfg,
		Log:     logger,
	}
	h.Router = mux.NewRouter()
	h.mapRoutes()

	h.Server = &http.Server{
		Addr:              address,
		Handler:           h.Router,
		ReadHeaderTimeout: 15 * time.Second,
	}

	return h
}

func (h *Handler) mapRoutes() {
	h.Router.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	m := NewMiddlewares(h.Log)

	fs := http.StripPrefix(h.Config.PodcastFeedsPathPrefix, http.FileServer(http.Dir(h.Service.GetFeedsRoot())))
	h.Router.PathPrefix(h.Config.PodcastFeedsPathPrefix).Handler(m.LoggingMiddleware(fs))
}

func (h *Handler) Serve() error {
	h.Log.Infoln("Starting server at", address)

	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			h.Log.WithError(err).Errorln("Server Error")
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := h.Server.Shutdown(ctx)
	if err != nil {
		h.Log.WithError(err).Errorln("Problem shutting down server")
		return err
	}

	h.Log.Infoln("Shut down server gracefully")

	return nil
}
