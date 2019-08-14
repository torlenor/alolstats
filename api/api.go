package api

import (
	"net/http"

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

func (a *API) run(listenAddress string) {
	a.log.Fatal(http.ListenAndServe(listenAddress, handlers.CORS()(a.router)))
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
	go a.run(listenAddress)
}

// Stop the REST API
func (a *API) Stop() {
	a.log.Println("Stopping REST API")
	// TODO
}
