package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"

	"git.abyle.org/hps/alolstats/config"
	"git.abyle.org/hps/alolstats/logging"
)

// API represents a Rest API instance of a ALoLStats instance
type API struct {
	config config.API
	router *mux.Router
	log    *logrus.Entry
	prefix string

	server *http.Server
	quit   chan interface{}
}

// NewAPI creates a new Rest API instance
func NewAPI(cfg config.API) (*API, error) {
	a := &API{
		config: cfg,
		router: mux.NewRouter(),
		log:    logging.Get("RestAPI"),
		prefix: "/v1",
	}

	return a, nil
}

// AttachModuleGet registers a new GET handler for the API
func (a *API) AttachModuleGet(path string, f func(http.ResponseWriter, *http.Request)) {
	a.log.Infoln("Registering GET handler:", a.prefix+path)
	a.router.HandleFunc(a.prefix+path, f).Methods("GET")
}

// AttachModulePost registers a new POST handler for the API
func (a *API) AttachModulePost(path string, f func(http.ResponseWriter, *http.Request)) {
	a.log.Infoln("Registering POST handler:", a.prefix+path)
	a.router.HandleFunc(a.prefix+path, f).Methods("POST")
}

func (a *API) run() {
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.log.Fatal("Could not start http server: %v\n", err)
	}
}

// Start the REST API
func (a *API) Start() {
	var listenAddress string
	if len(a.config.IP) > 0 && len(a.config.Port) > 0 {
		listenAddress = a.config.IP + ":" + a.config.Port
	} else if len(a.config.Port) > 0 {
		listenAddress = ":" + a.config.Port
	} else {
		a.log.Fatal("REST API activated but no valid configuration found. At least port has to specified in config file!")
	}

	a.log.Infof("Starting REST API on %s", listenAddress)

	a.server = &http.Server{
		Addr:    listenAddress,
		Handler: handlers.CORS()(a.router),
	}

	go a.run()
}

// Stop the REST API
func (a *API) Stop() {
	a.log.Println("Stopping REST API")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.server.SetKeepAlivesEnabled(false)
	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
}
