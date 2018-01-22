package stats

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"testing"
	"time"
)

var requiredMetrics = []string{promMetricNames[callsConst], promMetricNames[queuedConst],
	promMetricNames[completedConst], promMetricNames[failedConst], promMetricNames[runningConst], promMetricNames[timedoutConst], promMetricNames[errorsConst]}

// This file contains test utilities only (no tests)

func assertNoError(t *testing.T, assertionText string, err error) {
	if err != nil {
		t.Fatal(assertionText + " FAILED due to error: " + err.Error())
	}
}

func checkNoError(t *testing.T, assertionText string, err error) error {
	if err != nil {
		return errors.New(assertionText + " FAILED due to error: " + err.Error())
	}
	return nil
}

func assertIntsEqual(t *testing.T, assertionText string, expected int, actual int) {
	if actual != expected {
		t.Fatal(assertionText + " FAILED: expected " + strconv.Itoa(expected) + ", actual " + strconv.Itoa(actual))
	}
}

func checkIntsEqual(t *testing.T, assertionText string, expected int, actual int) error {
	if actual != expected {
		return errors.New(assertionText + " FAILED: expected " + strconv.Itoa(expected) + ", actual " + strconv.Itoa(actual))
	}
	return nil
}

func assertStringsEqual(t *testing.T, assertionText string, expected string, actual string) {
	if actual != expected {
		t.Fatal(assertionText + " FAILED: expected " + expected + ", actual " + actual)
	}
}

func checkStringsEqual(t *testing.T, assertionText string, expected string, actual string) error {
	if actual != expected {
		return errors.New(assertionText + " FAILED: expected " + expected + ", actual " + actual)
	}
	return nil
}

func assertNotNil(t *testing.T, assertionText string, actual interface{}) {
	if actual == nil {
		t.Fatal(assertionText + " FAILED: expected a non-nil value")
	}
}

func checkNotNil(t *testing.T, assertionText string, actual interface{}) error {
	if actual == nil {
		return errors.New(assertionText + " FAILED: expected a non-nil value")
	}
	return nil
}

// Return the sum of fn_completed for all applications and routes
// Metric values are obtained by scraping the /metrics endpoint directly
func getCompletedForAll(t *testing.T) int {
	return getMetricForAll(t, promMetricNames[completedConst])
}

// Return the sum of fn_completed for all routes in the specified application
// Metric values are obtained by scraping the /metrics endpoint directly
func getCompletedForApp(t *testing.T, appname string) int {
	return getMetricForApp(t, appname, promMetricNames[completedConst])
}

// Return the value of fn_completed for the specified application and route
// Metric values are obtained by scraping the /metrics endpoint directly
func getCompletedForAppAndRoute(t *testing.T, appname string, routename string) int {
	return getMetricForAppAndRoute(t, appname, routename, promMetricNames[completedConst])
}

// Return the sum of fn_failed for all applications and routes
// Metric values are obtained by scraping the /metrics endpoint directly
func getFailedForAll(t *testing.T) int {
	return getMetricForAll(t, promMetricNames[failedConst])
}

// Return the sum of fn_failed for all routes in the specified application
// Metric values are obtained by scraping the /metrics endpoint directly
func getFailedForApp(t *testing.T, appname string) int {
	return getMetricForApp(t, appname, promMetricNames[failedConst])
}

// Return the value of fn_failed for the specified application and route
// Metric values are obtained by scraping the /metrics endpoint directly
func getFailedForAppAndRoute(t *testing.T, appname string, routename string) int {
	return getMetricForAppAndRoute(t, appname, routename, promMetricNames[failedConst])
}

// Return the sum of fn_calls for all applications and routes
// Metric values are obtained by scraping the /metrics endpoint directly
func getCallsForAll(t *testing.T) int {
	return getMetricForAll(t, promMetricNames[callsConst])
}

// Return the sum of fn_calls for all routes in the specified application
// Metric values are obtained by scraping the /metrics endpoint directly
func getCallsForApp(t *testing.T, appname string) int {
	return getMetricForApp(t, appname, promMetricNames[callsConst])
}

// Return the value of fn_calls for the specified application and route
// Metric values are obtained by scraping the /metrics endpoint directly
func getCallsForAppAndRoute(t *testing.T, appname string, routename string) int {
	return getMetricForAppAndRoute(t, appname, routename, promMetricNames[callsConst])
}

// Return the sum of fn_errors for all applications and routes
// Metric values are obtained by scraping the /metrics endpoint directly
func getErrorsForAll(t *testing.T) int {
	return getMetricForAll(t, promMetricNames[errorsConst])
}

// Return the sum of fn_errors for all routes in the specified application
// Metric values are obtained by scraping the /metrics endpoint directly
func getErrorsForApp(t *testing.T, appname string) int {
	return getMetricForApp(t, appname, promMetricNames[errorsConst])
}

// Return the value of fn_errors for the specified application and route
// Metric values are obtained by scraping the /metrics endpoint directly
func getErrorsForAppAndRoute(t *testing.T, appname string, routename string) int {
	return getMetricForAppAndRoute(t, appname, routename, promMetricNames[errorsConst])
}

// Return the sum of fn_timedout for all applications and routes
// Metric values are obtained by scraping the /metrics endpoint directly
func getTimedOutForAll(t *testing.T) int {
	return getMetricForAll(t, promMetricNames[timedoutConst])
}

// Return the sum of fn_timedout for all routes in the specified application
// Metric values are obtained by scraping the /metrics endpoint directly
func getImedOutForApp(t *testing.T, appname string) int {
	return getMetricForApp(t, appname, promMetricNames[timedoutConst])
}

// Return the value of fn_timedout for the specified application and route
// Metric values are obtained by scraping the /metrics endpoint directly
func getTimedOutForAppAndRoute(t *testing.T, appname string, routename string) int {
	return getMetricForAppAndRoute(t, appname, routename, promMetricNames[timedoutConst])
}

// Return the sum of the specified Prometheus gauge or counter metric for all applications and routes
// Metric values are obtained by scraping the /metrics endpoint directly
func getMetricForAll(t *testing.T, metricname string) int {
	sum := 0
	for _, appname := range getApplications(t) {
		for _, routename := range getRoutes(t, appname) {
			metrics := getAllMetricsForAppAndRoute(t, appname, routename)
			sum += metrics[metricname]
		}
	}
	return sum
}

// Return the sums of all metrics for all applications and routes
// Metric values are obtained by scraping the /metrics endpoint directly
func getAllMetricsForAll(t *testing.T) map[string]int {
	result := make(map[string]int)
	for _, appname := range getApplications(t) {
		for _, routename := range getRoutes(t, appname) {
			allMetricsForAppAndRoute := getAllMetricsForAppAndRoute(t, appname, routename)
			for _, metricname := range requiredMetrics {
				result[metricname] = result[metricname] + allMetricsForAppAndRoute[metricname]
			}
		}
	}
	return result
}

// Return the sum of the specified Prometheus gauge or counter metric for all routes in the specified application
// Metric values are obtained by scraping the /metrics endpoint directly
func getMetricForApp(t *testing.T, appname string, metricname string) int {
	sum := 0
	for _, routename := range getRoutes(t, appname) {
		sum += getAllMetricsForAppAndRoute(t, appname, routename)[metricname]
	}
	return sum
}

// Return the sums of all metrics for all routes in the specified application
// Metric values are obtained by scraping the /metrics endpoint directly
func getAllMetricsForApp(t *testing.T, appname string) map[string]int {
	result := make(map[string]int)

	for _, routename := range getRoutes(t, appname) {
		allMetricsForAppAndRoute := getAllMetricsForAppAndRoute(t, appname, routename)
		for _, metricname := range requiredMetrics {
			result[metricname] = result[metricname] + allMetricsForAppAndRoute[metricname]
		}
	}
	return result
}

// Return the value of the specified Prometheus gauge or counter metric for the specified application and route
// Metric values are obtained by scraping the /metrics endpoint directly
func getMetricForAppAndRoute(t *testing.T, appname string, routename string, metricname string) int {
	return getAllMetricsForAppAndRoute(t, appname, routename)[metricname]
}

// Get all metrics for the specified application and name
// Metric values are obtained by scraping the /metrics endpoint directly
func getAllMetricsForAppAndRoute(t *testing.T, appname string, routename string) map[string]int {

	result := make(map[string]int)

	// get all Prometheus metrics
	scrapedMetrics := getURLAsString(t, "http://localhost:8080/metrics")

	for _, thisMetricName := range requiredMetrics {
		var thisMetricValue int
		var err error
		//regularExpression := thisMetricName + `{app="` + appname + `",path="/` + routename + `"} (\d+)`
		regularExpression := thisMetricName + `{` + appnameLabel + `="` + appname + `",` + pathLabel + `="/` + routename + `"} (\d+)`

		re := regexp.MustCompile(regularExpression)
		matches := re.FindStringSubmatch(scrapedMetrics)
		if len(matches) == 0 {
			thisMetricValue = 0
		} else {
			thisMetricValue, err = strconv.Atoi(matches[1])
			if err != nil {
				t.Fatal(err)
			}
		}
		result[thisMetricName] = thisMetricValue
	}
	return result

}

// Return the names of all applications
// This is obtained by using the /apps endpoint
func getApplications(t *testing.T) []string {
	responseAsJSON := getURLAsJSON(t, "http://localhost:8080/v1/apps")
	responseAsMap := responseAsJSON.(map[string]interface{})
	println(responseAsMap)
	if responseAsMap["next_cursor"].(string) != "" {
		t.Fatal(errors.New("Too many apps, need extend this code to use cursor"))
	}
	appsAsJSON := responseAsMap["apps"]
	appsAsArray := appsAsJSON.([]interface{})
	apps := make([]string, len(appsAsArray))
	for i, appDataAsJSON := range appsAsArray {
		appDataAsMap := appDataAsJSON.(map[string]interface{})
		apps[i] = appDataAsMap["name"].(string)
	}
	return apps
}

// Return the names of all routes in the specified application
// This is obtained by using the /apps/appname/routes endpoint
// The leading "/" is removed from the routes returned
func getRoutes(t *testing.T, appname string) []string {
	responseAsJSON := getURLAsJSON(t, "http://localhost:8080/v1/apps/"+appname+"/routes")
	responseAsMap := responseAsJSON.(map[string]interface{})
	println(responseAsMap)
	if responseAsMap["next_cursor"].(string) != "" {
		t.Fatal(errors.New("Too many routes, need extend this code to use cursor"))
	}
	routesAsJSON := responseAsMap["routes"]
	routesAsArray := routesAsJSON.([]interface{})
	routes := make([]string, len(routesAsArray))
	for i, routeDataAsJSON := range routesAsArray {
		routeDataAsMap := routeDataAsJSON.(map[string]interface{})
		routeName := routeDataAsMap["path"].(string)
		routes[i] = routeName[1:] // remove leading "/"
	}
	return routes
}

// GET the specified URL and return the result as a string
func getURLAsString(t *testing.T, url string) string {
	httpClient := http.Client{
		Timeout: time.Second * 60, // Maximum of 60 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	req.Header.Set("User-Agent", "github.com/fnproject/ext-statsapi/metrics_test")

	res, err := httpClient.Do(req)
	if err != nil {
		t.Fatal(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err.Error())
	}

	return string(body[:])
}

// GET the specified URL and return the result as unmarshalled JSON
func getURLAsJSON(t *testing.T, url string) interface{} {

	res, err := getURL(url)
	if err != nil {
		t.Fatal(err.Error())
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		t.Fatal(readErr.Error())
	}

	var m interface{}
	unmarshallErr := json.Unmarshal(body, &m)
	if unmarshallErr != nil {
		t.Fatal(readErr.Error())
	}
	return m
}

func getURL(url string) (*http.Response, error) {
	httpClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return &http.Response{}, err
	}

	req.Header.Set("User-Agent", "github.com/fnproject/ext-statsapi/stats-test")

	resp, getErr := httpClient.Do(req)
	if getErr != nil {
		return &http.Response{}, err
	}

	return resp, nil
}

func getURLWithPanic(url string) *http.Response {
	resp, err := getURL(url)
	if err != nil {
		panic(err)
	}
	return resp
}
