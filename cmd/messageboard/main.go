package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/guilherme-santos/messageboard"
	mbhttp "github.com/guilherme-santos/messageboard/http"
	"github.com/guilherme-santos/messageboard/mongodb"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	HTTPAddr          string
	Credentials       map[string]string
	MongoDBURL        string
	MongoDBInitialCSV string
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mgoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDBURL))
	if err != nil {
		log.Println("unable to connect to mongodb:", err)
		return
	}
	defer mgoClient.Disconnect(context.Background())

	storage := mongodb.NewMessageBoardStorage(mgoClient)

	log.Println("loading csv file", cfg.MongoDBInitialCSV)
	err = storage.LoadCSV(cfg.MongoDBInitialCSV)
	if err != nil {
		log.Println("unable to load csv file:", err)
		return
	}

	svc := messageboard.NewService(storage)

	// I'm using go-chi because it's lightweight (https://github.com/go-chi/chi#benchmarks) and simple
	// I usually reconfigure it, with nice logger and middlewares and so on,
	// but i want to keep it as simple as possible.
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	// Register message board handler to the router
	mbhttp.NewMessageBoardHandler(router, svc, cfg.Credentials)

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

	// Check if httpServer.ListenAndServe returned an error
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

	cfg.Credentials = make(map[string]string)

	creds := strings.Split(os.Getenv("CREDENTIALS"), ",")
	if len(creds) > 0 && creds[0] != "" {
		for i, cred := range creds {
			cred = strings.TrimSpace(cred)
			if cred == "" {
				continue
			}

			usrPasswd := strings.SplitN(cred, ":", 2)

			var user, passwd string
			if len(usrPasswd) == 2 {
				user = strings.TrimSpace(usrPasswd[0])
				passwd = strings.TrimSpace(usrPasswd[1])
			}
			if user == "" || passwd == "" {
				log.Printf("ignoring CREDENTIALS of position %d: %q", i, usrPasswd)
				continue
			}
			cfg.Credentials[usrPasswd[0]] = usrPasswd[1]
		}
	}

	cfg.MongoDBURL = os.Getenv("MONGODB_URL")
	cfg.MongoDBInitialCSV = os.Getenv("MONGODB_INITIAL_CSV")
	return nil
}
