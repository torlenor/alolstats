package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/torlenor/alolstats/api"
	"github.com/torlenor/alolstats/config"
	"github.com/torlenor/alolstats/logging"
	"github.com/torlenor/alolstats/memorybackend"
	"github.com/torlenor/alolstats/mongobackend"
	"github.com/torlenor/alolstats/riotclient"
	"github.com/torlenor/alolstats/statsrunner"
	"github.com/torlenor/alolstats/storage"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

const (
	defaultConfigPath   = "./alolstats.toml"
	defaultLoggingLevel = "debug"
)

/**
 * Version should be set while build using ldflags (see Makefile)
 */
var version string

var configPath string
var loggingLevel string

var interrupt chan os.Signal

var log *logrus.Entry

func init() {
	flag.StringVar(&configPath, "c", defaultConfigPath, "Path to toml config file")
	flag.StringVar(&loggingLevel, "l", defaultLoggingLevel, "Logging level (panic, fatal, error, warn/warning, info or debug)")
	flag.Parse()
}

func storageBackendCreator(cfg config.StorageBackend) (storage.Backend, error) {
	backendName := strings.ToLower(cfg.Backend)

	switch backendName {
	case "memory":
		backend, err := memorybackend.NewBackend()
		if err != nil {
			log.Errorln("Error creating the Memory Storage Backend:" + err.Error())
			return nil, err
		}
		return backend, nil
	case "mongo":
		fallthrough
	case "mongodb":
		backend, err := mongobackend.NewBackend(cfg.MongoBackend)
		if err != nil {
			log.Errorln("Error creating the Memory Storage Backend:" + err.Error())
			return nil, err
		}
		return backend, nil
	default:
		return nil, fmt.Errorf("Unknown storage backend specified in config: %s", cfg.Backend)
	}
}

func main() {
	logging.Init()
	logging.SetLoggingLevel(loggingLevel)

	log = logging.Get("main")

	log.Println("ALoLStats (" + version + ") is STARTING")

	interrupt = make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var cfg config.Config
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		log.Fatalln(err)
	}

	api, err := api.NewAPI(cfg.API)
	if err != nil {
		log.Fatalln("Error creating the API:" + err.Error())
	}

	client, err := riotclient.NewClient(&http.Client{}, cfg.RiotClient)
	if err != nil {
		log.Fatalln("Error creating the Riot Client:" + err.Error())
	}
	client.Start()

	backend, err := storageBackendCreator(cfg.StorageBackend)
	if err != nil {
		log.Fatalln("Error creating the Storage Backend:" + err.Error())
	}

	storage, err := storage.NewStorage(cfg.LoLStorage, client, backend)
	if err != nil {
		log.Fatalln("Error creating the Storage:" + err.Error())
	}
	storage.RegisterAPI(api)
	storage.Start()

	statsRunner, err := statsrunner.NewStatsRunner(storage)
	if err != nil {
		log.Fatalln("Error creating the StatsRunner:" + err.Error())
	}
	statsRunner.RegisterAPI(api)

	api.Start()

	log.Println("ALoLStats (" + version + ") is READY")

	for {
		select {
		case <-interrupt:
			api.Stop()
			storage.Stop()
			client.Stop()
			log.Printf("Storage handeled %d requests since startup", storage.GetHandeledRequests())
			log.Println("ALoLStats gracefully shut down")
			return
		}
	}
}
