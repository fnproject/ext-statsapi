package statistics

import (
	"encoding/json"
	"fmt"
	"github.com/fnproject/ext-metrics/fncommon"
	"github.com/fnproject/fn/api/models"
	"github.com/fnproject/fn/api/server"
	"net/http"
)

const (
	EnvPromHost = "fn_ext_metrics_prom_host"
	EnvPromPort = "fn_ext_metrics_prom_port"
)

var promHost string
var promPort string

type globalStatisticsHandler struct{}
type appStatisticsHandler struct{}
type routeStatisticsHandler struct{}

func init() {
	promHost = fncommon.GetEnv(EnvPromHost, "localhost")
	promPort = fncommon.GetEnv(EnvPromPort, "9090")
}

// Call this function from the main function of any Fn server to add this extension to the API
func AddEndpoints(s *server.Server) {

	s.AddEndpoint("GET", "/statistics", &globalStatisticsHandler{})

	// the following will be at /v1/apps/:app_name/statistics
	s.AddAppEndpoint("GET", "/statistics", &appStatisticsHandler{})

	// the following will be at /v1/apps/:app_name/routes/:route_name/statistics
	s.AddRouteEndpoint("GET", "/statistics", &routeStatisticsHandler{})
}

func (h *globalStatisticsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jsonData := handle(r, "", "")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonData))
}

func (h *appStatisticsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, app *models.App) {
	jsonData := handle(r, app.Name, "")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonData))
}

func (h *routeStatisticsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, app *models.App, route *models.Route) {
	jsonData := handle(r, app.Name, route.Path)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonData))
}

// query scopes
const (
	query_scope_global = iota
	query_scope_app    = iota
	query_scope_route  = iota
)

// metric types
const (
	completed = iota
	failed    = iota
	durations = iota
)

// keys in the returned JSON data array, keyed by metric type
var jsonKeys = map[int]string{
	completed: "completed",
	failed:    "failed",
	durations: "durations",
}

// Process the request and return the requested data as JSON
func handle(r *http.Request, appName string, routeName string) []byte {

	// parse query params and provide default values if needed
	startTimeString, endTimeString, stepString, err := getQueryParams(r)
	if err != nil {
		return getErrorAsJSON(err)
	}

	var queryScope int
	if appName == "" {
		queryScope = query_scope_global
	} else if routeName == "" {
		queryScope = query_scope_app
	} else {
		queryScope = query_scope_route
	}

	// create a struct that will contain our response prior to conversion to JSONs
	responseStruct := new(metricsResponse)
	responseStruct.Status = "success"
	responseStruct.Data = make(map[string][]metricsTimeValuePair)

	// for each metric type, query Prometheus and populate the response struct
	for metricType, jsonKey := range jsonKeys {
		// construct the Prometheus request URL
		url := buildPrometheusRequest(queryBuilders[metricType], promHost, promPort, queryScope, metricType, appName, routeName, startTimeString, endTimeString, stepString)
		// execute the Prometheus request and extract the array of time-value pairs from the response
		metricDataArray, err := executePrometheusRequest(url)
		if err != nil {
			return getErrorAsJSON(err)
		}
		responseStruct.Data[jsonKey] = metricDataArray
	}

	// convert the response struct to JSON
	jsonData, err := json.Marshal(responseStruct)
	if err != nil {
		return getErrorAsJSON(err)
	}
	return jsonData
}

func getErrorAsJSON(err error) []byte {
	errorResponseStruc := new(errorResponse)
	errorResponseStruc.Status = STATS_STATUS_ERROR
	errorResponseStruc.Error = err.Error()
	errorJsonData, _ := json.Marshal(errorResponseStruc)
	return errorJsonData
}
