package stats

import (
	"encoding/json"
	"fmt"
	"github.com/fnproject/ext-statsapi/fncommon"
	"github.com/fnproject/fn/api/models"
	"github.com/fnproject/fn/fnext"
	"net/http"
)

const (
	EnvPromHost = "FN_EXT_STATS_PROM_HOST"
	EnvPromPort = "FN_EXT_STATS_PROM_PORT"
)

var promHost string
var promPort string

type globalStatisticsHandler struct{}
type appStatisticsHandler struct{}
type routeStatisticsHandler struct{}

func AddEndpoints(s fnext.ExtServer) {

	promHost = fncommon.GetEnv(EnvPromHost, "localhost")
	promPort = fncommon.GetEnv(EnvPromPort, "9090")

	s.AddEndpoint("GET", "/stats", &globalStatisticsHandler{})
	s.AddEndpoint("GET", "/statistics", &globalStatisticsHandler{})

	// the following will be at /v1/apps/:app_name/stats
	s.AddAppEndpoint("GET", "/stats", &appStatisticsHandler{})
	s.AddAppEndpoint("GET", "/statistics", &appStatisticsHandler{})

	// the following will be at /v1/apps/:app_name/routes/:route_name/stats
	s.AddRouteEndpoint("GET", "/stats", &routeStatisticsHandler{})
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

// these constants represent the various types of statistic returned by this API
// If you add a new type of statistic you must also
// (0) add a new entry to the map jsonKeys below
// (1) add new entries to the maps promMetricNames in build_prometheus_request.go
// (2) update the map queryBuilders (in build_prometheus_request.go)
//     with the name of the appropriate query builder for this metric type
//     this is essentially a case of specifying whether the metric is a histogram or a counter/gauge
const (
	completedConst = iota
	failedConst    = iota
	durationsConst = iota
	callsConst     = iota
	errorsConst    = iota
	timedoutConst  = iota
)

// in this map, the key is the constant for the type of statistic and
// the corresponding value is the name of the key that will hold this type of statistic in the returned JSON data structure
// see comment above for information on adding a new type of statistic
var jsonKeys = map[int]string{
	completedConst: "completed",
	failedConst:    "failed",
	durationsConst: "durations",
	callsConst:     "calls",
	errorsConst:    "errors",
	timedoutConst:  "timedout",
}

var appLabel = "fn_appname"
var routeLabel = "fn_path"

// Process the request and return the requested data as JSON
func handle(r *http.Request, appName string, routeName string) []byte {

	// parse query params and provide default values if needed
	startTimeString, endTimeString, stepString, err := getQueryParams(r)
	if err != nil {
		return getErrorAsJSON(err)
	}

	// create a struct that will contain our response prior to conversion to JSONs
	responseStruct := new(metricsResponse)
	responseStruct.Status = "success"
	responseStruct.Data = make(map[string][]metricsTimeValuePair)

	// for each metric type, query Prometheus and populate the response struct
	for metricType, jsonKey := range jsonKeys {
		// construct the Prometheus request URL
		url := buildPrometheusRequest(queryBuilders[metricType], promHost, promPort, metricType, appName, routeName, startTimeString, endTimeString, stepString)
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
