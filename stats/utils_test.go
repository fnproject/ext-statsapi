package stats

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"
)

// This file contains test utilities only (no tests)

func assertNoError(t *testing.T, assertionText string, err error) {
	if err != nil {
		t.Fatal(assertionText + " FAILED due to error: " + err.Error())
	}
}

func assertIntsEqual(t *testing.T, assertionText string, expected int, actual int) {
	if actual != expected {
		t.Fatal(assertionText + " FAILED: expected " + strconv.Itoa(expected) + ", actual " + strconv.Itoa(actual))
	}
}

func assertStringsEqual(t *testing.T, assertionText string, expected string, actual string) {
	if actual != expected {
		t.Fatal(assertionText + " FAILED: expected " + expected + ", actual " + actual)
	}
}

func assertNotNil(t *testing.T, assertionText string, actual interface{}) {
	if actual == nil {
		t.Fatal(assertionText + " FAILED: expected a non-nil value")
	}
}

// Return the sum of fn_api_completed for all applications and routes
// Metric values are obtained by scraping the /metrics endpoint directly
func getCompleted(t *testing.T) int {
	completed := 0
	for _, appname := range getApplications(t) {
		for _, routename := range getRoutes(t, appname) {
			metrics := getMetrics(t, appname, routename)
			completed += metrics["fn_api_completed"]
		}
	}
	return completed
}

// Return the sum of fn_api_completed for all routes in the specified application
// Metric values are obtained by scraping the /metrics endpoint directly
func getCompletedForApp(t *testing.T, appname string) int {
	completed := 0
	for _, routename := range getRoutes(t, appname) {
		completed += getMetrics(t, appname, routename)["fn_api_completed"]
	}
	return completed
}

// Return the value of fn_api_completed the specified application and route
// Metric values are obtained by scraping the /metrics endpoint directly
func getCompletedForAppAndRoute(t *testing.T, appname string, routename string) int {
	return getMetrics(t, appname, routename)["fn_api_completed"]
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
	httpClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	req.Header.Set("User-Agent", "github.com/fnproject/ext-statsapi/main-test")

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		t.Fatal(getErr.Error())
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
