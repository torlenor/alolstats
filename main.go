package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/torlenor/alolstats/api"
	"github.com/torlenor/alolstats/config"
	"github.com/torlenor/alolstats/logging"
	"github.com/torlenor/alolstats/memorybackend"
	"github.com/torlenor/alolstats/riotclient"
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

func main() {
	logging.Init()
	logging.SetLoggingLevel(loggingLevel)

	log = logging.Get("main")

	log.Println("ALoLStats (" + version + ") is STARTING")

	interrupt = make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var cfg config.Config
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		fmt.Println(err)
		return
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

	backend, err := memorybackend.NewBackend()
	if err != nil {
		log.Fatalln("Error creating the Memory Storage Backend:" + err.Error())
	}

	storage, err := storage.NewStorage(cfg.LoLStorage, client, backend)
	if err != nil {
		log.Fatalln("Error creating the Storage:" + err.Error())
	}

	storage.RegisterAPI(api)

	storage.Start()
	api.Start()

	log.Println("ALoLStats (" + version + ") is READY")

	for {
		select {
		case <-interrupt:
			api.Stop()
			storage.Stop()
			log.Printf("Storage handeled %d requests since startup", storage.GetHandeledRequests())
			log.Println("ALoLStats gracefully shut down")
			return
		}
	}
}
