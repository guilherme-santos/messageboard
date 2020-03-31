package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/guilherme-santos/messageboard"
	mbhttp "github.com/guilherme-santos/messageboard/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Config struct {
	HTTPAddr string
}

var cfg Config

func init() {
	// I usually use https://github.com/kelseyhightower/envconfig
	// but for sake of simplicity, I'll do by hand.
	err := loadConfig(&cfg)
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	var storage messageboard.Storage

	svc := messageboard.NewService(storage)

	// I'm using go-chi because it's lightweight (https://github.com/go-chi/chi#benchmarks) and simple
	// I usually reconfigure it, with nice logger and middlewares and so on,
	// but i want to keep it as simple as possible.
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	// Register message board handler to the router
	mbhttp.NewMessageBoardHandler(router, svc)

	httpServer := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}

	log.Println("running webserver on", httpServer.Addr)

	// Run http server in a goroutine to be able to catch SIGINT
	errCh := make(chan error)
	go func() {
		defer close(errCh)
		err := httpServer.ListenAndServe()
		errCh <- err
	}()

	select {
	case err := <-errCh:
		log.Println("unable to run webserver:", err)
		return
	case <-time.After(time.Second):
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	// Wait until receive one of the signals
	<-sigCh

	log.Println("signal received, shutting down webserver")
	httpServer.Shutdown(context.Background())
}

func loadConfig(cfg *Config) error {
	cfg.HTTPAddr = os.Getenv("HTTP_ADDR")
	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = "0.0.0.0:80"
	}
	return nil
}
