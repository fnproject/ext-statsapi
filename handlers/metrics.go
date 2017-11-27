package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fnproject/fn/api/models"
	"github.com/fnproject/fn/api/server"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type globalStatisticsHandler struct{}
type appStatisticsHandler struct{}

// Call this function from the main function of any Fn server to add this extension to the API
func AddEndpoints(s *server.Server) {
	s.AddEndpoint("GET", "/statistics", &globalStatisticsHandler{})

	// the following will be at /v1/apps/:app_name/statistics
	s.AddAppEndpoint("GET", "/statistics", &appStatisticsHandler{})
}

func (h *globalStatisticsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "globalStatisticsHandler handling %q", html.EscapeString(r.URL.Path))
	jsonData := getJSONResponse(r, "", "")
    w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonData))
}

func (h *appStatisticsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, app *models.App) {
	//fmt.Fprintf(w, "appStatisticsHandler handling %q", html.EscapeString(r.URL.Path))
	jsonData := getJSONResponse(r, app.Name, "")
    w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonData))
}

// get the response as JSON
func getJSONResponse(r *http.Request, appname string, routename string) []byte {
	// TODO common func to handle errors

	// parse query params and provide default values if needed
	starttimeString, endtimeString, step, err := getParams(r)
	if err != nil {
		errorResponseStruc := new(errorResponse)
		errorResponseStruc.Status = STATUS_ERROR
		errorResponseStruc.Error = err.Error()
		errorJsonData, _ := json.Marshal(errorResponseStruc)
		return errorJsonData
	}

	// get response as a struct
	responseObject, err := getMetricsResponse(appname, routename, starttimeString, endtimeString, step)
	if err != nil {
		errorResponseStruc := new(errorResponse)
		errorResponseStruc.Status = STATUS_ERROR
		errorResponseStruc.Error = err.Error()
		errorJsonData, _ := json.Marshal(errorResponseStruc)
		return errorJsonData
	}

	jsonData, err := json.Marshal(responseObject)
	if err != nil {
		errorResponseStruc := new(errorResponse)
		errorResponseStruc.Status = STATUS_ERROR
		errorResponseStruc.Error = err.Error()
		errorJsonData, _ := json.Marshal(errorResponseStruc)
		return errorJsonData
	}
	return jsonData
}

// Extract and return the required URL query parameters, generating default values if missing

// starttimeString - start of range - (string of format 2006-01-02T15:04:05.999Z07:00)
// endtimeString - end of range - (string of format 2006-01-02T15:04:05.999Z07:00)
// stepTime - milliseconds - time duration between values
//
func getParams(r *http.Request) (string, string, string, error) {

	var starttimeString, endtimeString, stepString string
	var starttime, endtime time.Time
	var err error

	startTimeParams := r.URL.Query()["starttime"]
	if len(startTimeParams) > 0 {
		starttimeString = startTimeParams[0]
		starttime, err = time.Parse(prometheusTimeFormat, starttimeString)
		if err != nil {
			return "", "", "", errors.New("Unable to parse starttime parameter: " + err.Error())
		}
	}

	endTimeParams := r.URL.Query()["endtime"]
	if len(endTimeParams) > 0 {
		endtimeString = endTimeParams[0]
		endtime, err = time.Parse(prometheusTimeFormat, endtimeString)
		if err != nil {
			return "", "", "", errors.New("Unable to parse endtime parameter: " + err.Error())
		}
	}

	stepParams := r.URL.Query()["step"]
	if len(stepParams) > 0 {
		stepString = stepParams[0]
		_, err = time.ParseDuration(stepString)
		if err != nil {
			return "", "", "", errors.New("Unable to parse step parameter: " + err.Error())
		}
	}

	switch {
	case len(startTimeParams) == 0 && len(endTimeParams) == 0:
		// neither starttime or endtime specified, set starttime to 5m ago, set endtime to now
		endtime = time.Now()
		endtimeString = endtime.Format(prometheusTimeFormat)
		starttime = endtime.Add(-(time.Duration(5) * time.Minute))
		starttimeString = starttime.Format(prometheusTimeFormat)
		//println("getParams: neither starttime or endtime specified - set default starttime=", starttimeString, ", set default endtime=", endtimeString)
	case len(startTimeParams) == 0:
		// endtime is specified, starttime is not specified, set to 5mins before endtime
		starttime = endtime.Add(-(time.Duration(5) * time.Minute))
		starttimeString = starttime.Format(prometheusTimeFormat)
		//println("getParams: starttime not specified - set default starttime=", starttimeString, ", endtime=", endtimeString)
	case len(endTimeParams) == 0:
		// starttime is specified, endtime is not specified, set to now
		endtime = time.Now()
		endtimeString = endtime.Format(prometheusTimeFormat)
		//println("getParams: endtime not specified - starttime=", starttimeString, ", set default endtime=", endtimeString)
	default:
		// both starttime and endtime specified
	}

	if endtime.Before(starttime) {
		println("endtime (" + endtimeString + ") is before starttime (" + starttimeString + ")")
		err = errors.New("endtime (" + endtimeString + ") is before starttime (" + starttimeString + ")")
		return "", "", "", err
	}

	if len(stepParams) == 0 {
		// step not specified, assume 30 secs
		stepString = time.Duration(30 * time.Second).String()
		println("getParams: step not specified. step=", stepString)
	}

	return starttimeString, endtimeString, stepString, err
}

// the following structs represent the JSON returned by calls to the Prometheus API at api/v1/query_range

const STATUS_ERROR string = "error"
const STATUS_SUCCESS string = "success"

type promQueryRangeData struct {
	Status    string `json:"status"`    // "success" | "error" | ?
	Error     string `json:"errorType"` // only present if status=error
	ErrorType string `json:"error"`     // only present if status=error
	Data      data   `json:"data"`
}

type data struct {
	ResultType string         `json:"resultType"` // "matrix" (Range vectors) | "vector" (Instant vectors) | "scalar" (Scalar results) | "string" (String results)
	Result     []matrixResult `json:"result"`
}

type matrixResult struct { // Used when resultType is matrix
	Metric map[string]string `json:"metric"` // Map of label_name to label_value (empty if this is a sum)
	Value  []timeValuePair   `json:"values"`
}

type timeValuePair []interface{} // This is an array with two elements, <unix_time>, "<scalar_value>"

func (tvp timeValuePair) UnixTime() float64 {
	result, _ := tvp[0].(float64)
	return result
}

func (tvp timeValuePair) ScalarValue() string {
	result, _ := tvp[1].(string)
	return result
}

// time format required when sending queries to Prometheus
const prometheusTimeFormat = "2006-01-02T15:04:05.999Z07:00"

func getMetricsResponse(appName string, routeName string, starttimeString string, endtimeString string, stepString string) (*metricsResponse, error) {

	responseStruct := new(metricsResponse)
	responseStruct.Status = "success"

	// make the map (key=metricName, value= array of metricsTimeValuePair strucs)
	responseStruct.Data = make(map[string][]metricsTimeValuePair)

	var requiredMetrics = map[string]string{
		"completed": "fn_api_completed",
		"failed":    "fn_api_failed",
		//"duration": "rate(fn_span_agent_cold_exec_duration_seconds_sum[1m])/rate(fn_span_agent_cold_exec_duration_seconds_count[1m])",
	}
	for jsonKey, metricName := range requiredMetrics {
		metricDataArray, err := getDataFor(metricName, appName, routeName, starttimeString, endtimeString, stepString)
		if err != nil {
			return nil, err
		}
		responseStruct.Data[jsonKey] = metricDataArray
	}
	return responseStruct, nil
}

// Get the data for a particular metric
func getDataFor(metricName string, appName string, routeName string, starttimeString string, endtimeString string, stepString string) ([]metricsTimeValuePair, error) {

	url := "http://localhost:9090/api/v1/query_range?query=sum(" + metricName + ")&start=" + starttimeString + "&end=" + endtimeString + "&step=" + stepString
	promClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	fmt.Println("== URL sent to Prometheus= ===========================================")
	fmt.Println(url)
	fmt.Println("===End of URL sent to Prometheus =====================================")

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "github.com/fnproject/ext-metrics")

	res, doErr := promClient.Do(req)
	if doErr != nil {
		return nil, doErr
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	fmt.Println("== JSON returned from Prometheus= ====================================")
	fmt.Println(string(body[:]))
	fmt.Println("===End of JSON returned from Prometheus ==============================")

	// Assume result is of type "matrix" (meaning this is a range vector)

	thisPromQueryRangeData := promQueryRangeData{}
	jsonErr := json.Unmarshal(body, &thisPromQueryRangeData)
	if jsonErr != nil {
		return nil, jsonErr
	}

	fmt.Println(thisPromQueryRangeData)

	if thisPromQueryRangeData.Status != "success" {
		return nil, errors.New("Error from Prometheus: " + thisPromQueryRangeData.ErrorType + ": " + thisPromQueryRangeData.Error)
	}

	if len(thisPromQueryRangeData.Data.Result) > 1 {
		//  we must have got the query wrong
		return nil, errors.New("data array returned by Prometheus has more than one element")
	}

	// how many data time-value pairs have we been given
	var numberOfTimeValuePairs int
	if len(thisPromQueryRangeData.Data.Result) == 0 {
		numberOfTimeValuePairs = 0
	} else {
		numberOfTimeValuePairs = len(thisPromQueryRangeData.Data.Result[0].Value)
	}

	// make the array for this metricName
	metricDataArray := make([]metricsTimeValuePair, numberOfTimeValuePairs)

	if numberOfTimeValuePairs > 0 {
		// populate the array with zero or more metricsTimeValuePair strucs
		for i, val := range thisPromQueryRangeData.Data.Result[0].Value {
			tvp := new(metricsTimeValuePair)
			tvp.Time = int64(val.UnixTime())
//			if val.ScalarValue()=="NaN" {
//				println("foo")
//			}
			value, err := strconv.ParseFloat(val.ScalarValue(), 64)
			tvp.Value = value
			if err != nil {
				return nil, errors.New("Error converting " + val.ScalarValue() + " to a float64")
			}
			metricDataArray[i] = *tvp
		}
	}

	return metricDataArray, nil
}

// the following structs represent the JSON returned by calls to the extension API
type errorResponse struct {
	Status string `json:"status"` // "error" (STATUS_ERROR)
	Error  string `json:"error"`  //  if Status is "error", set to the error message
}

// TODO have different structs for error and success
type metricsResponse struct {
	Status string                            `json:"status"` // "success" (STATUS_SUCCESS)
	Data   map[string][]metricsTimeValuePair `json:"data"`
}

type metricsTimeValuePair struct {
	Time  int64
	Value float64
}
