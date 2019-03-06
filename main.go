package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/torlenor/alolstats/api"
	"github.com/torlenor/alolstats/config"
	"github.com/torlenor/alolstats/fetchrunner"
	"github.com/torlenor/alolstats/logging"
	"github.com/torlenor/alolstats/memorybackend"
	"github.com/torlenor/alolstats/mongobackend"
	"github.com/torlenor/alolstats/riotclient"
	"github.com/torlenor/alolstats/riotclientv4"
	"github.com/torlenor/alolstats/statsrunner"
	"github.com/torlenor/alolstats/storage"

	"github.com/torlenor/alolstats/riotclient/datadragon"
	"github.com/torlenor/alolstats/riotclient/ratelimit"

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
var compTime string

var configPath string
var loggingLevel string
var logFile string

var interrupt chan os.Signal

var log *logrus.Entry

func init() {
	flag.StringVar(&configPath, "c", defaultConfigPath, "Path to toml config file")
	flag.StringVar(&loggingLevel, "l", defaultLoggingLevel, "Logging level (panic, fatal, error, warn/warning, info or debug)")
	flag.StringVar(&logFile, "L", "", "Log file to use")
	flag.Parse()
}

func statusEndpoint(w http.ResponseWriter, r *http.Request) {

	status := fmt.Sprintf(`{"status":"OK", "version":"%s", "compiled_at":"%s"}`, version, compTime)

	io.WriteString(w, string(status))
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

func riotClientCreator(cfg config.RiotClient) (riotclient.Client, error) {
	version := strings.ToLower(cfg.APIVersion)

	switch version {
	case "v3":
		return nil, fmt.Errorf("API v3 is not supported anymore")
	case "v4":
		httpClient := &http.Client{}
		ddragon, err := riotclientdd.New(httpClient, cfg)
		if err != nil {
			log.Errorln("Error creating Riot Client Data Dragon:" + err.Error())
			return nil, err
		}
		rateLimit, err := riotclientrl.New()
		if err != nil {
			log.Errorln("Error creating Riot Client Rate Limit Checker:" + err.Error())
			return nil, err
		}

		riotClient, err := riotclientv4.NewClient(httpClient, cfg, ddragon, rateLimit)
		if err != nil {
			log.Errorln("Error creating RiotClient APIVersion V4:" + err.Error())
			return nil, err
		}

		return riotClient, nil
	default:
		return nil, fmt.Errorf("Unknown RiotClient APIVersion specified in config: %s", cfg.APIVersion)
	}
}

func main() {
	logging.Init()

	if len(logFile) > 0 {
		logging.SetLogFile(logFile)
	}

	logging.SetLoggingLevel(loggingLevel)

	log = logging.Get("main")

	log.Println("ALoLStats " + version + " (" + compTime + ") is STARTING")

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
	api.AttachModuleGet("/status", statusEndpoint)

	clients := make(map[string]riotclient.Client)
	for name, clientConfig := range cfg.RiotClient {
		client, err := riotClientCreator(clientConfig)
		if err != nil {
			log.Fatalf("Error creating the Riot Client %s: %s", name, err.Error())
		}
		client.Start()
		clients[name] = client
	}

	backend, err := storageBackendCreator(cfg.StorageBackend)
	if err != nil {
		log.Fatalln("Error creating the Storage Backend:" + err.Error())
	}

	storage, err := storage.NewStorage(cfg.LoLStorage, clients, backend)
	if err != nil {
		log.Fatalln("Error creating the Storage:" + err.Error())
	}
	storage.RegisterAPI(api)
	storage.Start()

	statsRunner, err := statsrunner.NewStatsRunner(cfg.StatsRunner, storage)
	if err != nil {
		log.Fatalln("Error creating the StatsRunner:" + err.Error())
	}
	statsRunner.RegisterAPI(api)

	var fetchRunners []*fetchrunner.FetchRunner
	for name, fetchRunnerConfig := range cfg.FetchRunner {
		fetchRunner, err := fetchrunner.NewFetchRunner(fetchRunnerConfig, storage)
		if err != nil {
			log.Fatalf("Error creating the FetchRunner %s: %s", name, err.Error())
		}
		fetchRunner.Start()
		fetchRunners = append(fetchRunners, fetchRunner)
	}

	statsRunner.Start()
	api.Start()

	log.Println("ALoLStats (" + version + ") is READY")

	for {
		select {
		case <-interrupt:
			api.Stop()
			statsRunner.Stop()
			for _, fetchRunner := range fetchRunners {
				fetchRunner.Stop()
			}
			storage.Stop()
			for _, client := range clients {
				client.Stop()
			}
			log.Printf("Storage handeled %d requests since startup", storage.GetHandeledRequests())
			log.Printf("StatsRunner handeled %d requests since startup", statsRunner.GetHandeledRequests())
			log.Println("ALoLStats gracefully shut down")
			return
		}
	}
}
