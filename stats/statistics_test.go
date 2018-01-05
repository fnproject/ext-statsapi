package stats

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

// Requires the following to be running
// - an extended Fn server (such as the one in this project)
// - a Prometheus server, configured to scrape data from the Fn server
// Also requires the following to be run beforehand
//   bash test/create.bash
//   bash test/run-cold-sync.bash
//   bash test/run-cold-async.bash
//   bash test/run-hot-sync.bash
//   bash test/run-hot-async.bash

// Test a query which will return an error from the extension (rather than from the main router, such as app or route not found)
func TestBadRoute(t *testing.T) {
	url := "http://localhost:8080/v1/apps/hello-cold-sync-a/routes/hello-cold-sync-a1/stats?step=Wombat"
	response := getJSON(t, url)
	verifyFailedJSON(t, response, "Unable to parse step parameter: time: invalid duration Wombat")
}

// Test sync cold
func TestNeverCalled(t *testing.T) {
	// verify stats for a function that has never been called (since server startup)
	url := "http://localhost:8080/v1/apps/hello-cold-async-b/routes/hello-cold-async-b3/stats"
	response := getJSON(t, url)
	verifySuccessfulJSON(t, response, 0)
}

func TestAllFuncs(t *testing.T) {
	// verify stats for all functions
	// Assumes all the following have been run
	// test/run-cold-sync.bash
	// test/run-cold-async.bash
	// test/run-hot-sync.bash
	// test/run-hot-async.bash
	url := "http://localhost:8080/v1/stats"
	response := getJSON(t, url)
	verifySuccessfulJSON(t, response, 120)
}

func TestAllFuncsPerApp(t *testing.T) {
	// verify stats across all three functions in the app hello-cold-async-a
	// Assumes all the following have been run
	// test/run-cold-sync.bash
	// test/run-cold-async.bash
	url := "http://localhost:8080/v1/apps/hello-cold-async-a/stats"
	response := getJSON(t, url)
	verifySuccessfulJSON(t, response, 60)
}

// Test sync cold
// Assumes test/run-cold-sync.bash has been run
func TestSyncCold(t *testing.T) {
	// verify stats for sync cold functions
	url := "http://localhost:8080/v1/apps/hello-cold-sync-a/routes/hello-cold-sync-a1/stats"
	response := getJSON(t, url)
	verifySuccessfulJSON(t, response, 10)
}

// Test async cold
// Assumes test/run-cold-async.bash has been run
func TestAsyncCold(t *testing.T) {
	// verify stats for async cold functions
	url := "http://localhost:8080/v1/apps/hello-cold-async-a/routes/hello-cold-async-a2/stats"
	response := getJSON(t, url)
	verifySuccessfulJSON(t, response, 20)
}

// Test sync hot
// Assumes test/run-hot-sync.bash has been run
func TestSyncHot(t *testing.T) {
	// verify stats for sync hot functions
	url := "http://localhost:8080/v1/apps/hello-hot-sync-a/routes/hello-hot-sync-a1/stats"
	response := getJSON(t, url)
	verifySuccessfulJSON(t, response, 10)
}

// Test async hot
// Assumes test/run-hot-async.bash has been run
func TestAsyncHot(t *testing.T) {
	// verify stats for async hot functions
	url := "http://localhost:8080/v1/apps/hello-hot-async-a/routes/hello-hot-async-a1/stats"
	response := getJSON(t, url)
	verifySuccessfulJSON(t, response, 20)
}

func getJSON(t *testing.T, url string) interface{} {
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

// Verify the JSON that is returned from a successful call
func verifySuccessfulJSON(t *testing.T, response interface{}, expectedCompleted int) {

	responseAsMap := response.(map[string]interface{})

	// top level of JSON should be two keys, status and data
	assertIntsEqual(t, "Number of keys in top-level JSON", 2, len(responseAsMap))

	// check status
	assertNotNil(t, "status field should be present", responseAsMap["status"])
	statusAsString := responseAsMap["status"].(string)
	assertStringsEqual(t, "Status field", "success", statusAsString)

	// check data
	assertNotNil(t, "data field should be present", responseAsMap["data"])
	dataAsMap := responseAsMap["data"].(map[string]interface{})
	assertIntsEqual(t, "Number of keys in data", 3, len(dataAsMap))

	// check data -> completed
	assertNotNil(t, "completed field should be present", dataAsMap["completed"])
	completedAsArray := dataAsMap["completed"].([]interface{})
	// if any time-value pairs are present, verify them
	for _, timeAndValuePair := range completedAsArray {
		timeAndValuePairAsMap := timeAndValuePair.(map[string]interface{})
		// check data -> completed -> time
		assertNotNil(t, "time field should be present", timeAndValuePairAsMap["time"])
		// check data -> completed -> value
		assertNotNil(t, "value field should be present", timeAndValuePairAsMap["value"])
	}
	// verify latest completed count matches expectedCompleted
	if len(completedAsArray) > 0 {
		lastTimeAndValuePair := completedAsArray[len(completedAsArray)-1]
		lastTimeAndValuePairAsMap := lastTimeAndValuePair.(map[string]interface{})
		lastValue := lastTimeAndValuePairAsMap["value"]
		lastValueAsFloat64 := lastValue.(float64)
		lastvalueAsInt := int(lastValueAsFloat64)
		assertIntsEqual(t, "Completed count", expectedCompleted, lastvalueAsInt)
	} else {
		// Completed array is empty which implies a completed count of zero
		assertIntsEqual(t, "Completed array is empty", expectedCompleted, 0)
	}

	// check data -> failed
	assertNotNil(t, "failed field should be present", dataAsMap["failed"])
	failedAsArray := dataAsMap["completed"].([]interface{})
	// if any time-value pairs are present, verify them
	for _, timeAndValuePair := range failedAsArray {
		timeAndValuePairAsMap := timeAndValuePair.(map[string]interface{})

		// check data -> completed -> time
		assertNotNil(t, "time field should be present", timeAndValuePairAsMap["time"])

		// check data -> completed -> value
		assertNotNil(t, "value field should be present", timeAndValuePairAsMap["value"])
	}

	// check data -> durations
	assertNotNil(t, "durations field should be present", dataAsMap["durations"])
	durationsAsArray := dataAsMap["completed"].([]interface{})
	// if any time-value pairs are present, verify them
	for _, timeAndValuePair := range durationsAsArray {
		timeAndValuePairAsMap := timeAndValuePair.(map[string]interface{})

		// check data -> completed -> time
		assertNotNil(t, "time field should be present", timeAndValuePairAsMap["time"])

		// check data -> completed -> value
		assertNotNil(t, "value field should be present", timeAndValuePairAsMap["value"])
	}
}

// Verify the basic shape of the JSON that is returned for a failed call
func verifyFailedJSON(t *testing.T, response interface{}, expectedError string) {

	responseAsMap := response.(map[string]interface{})

	// top level of JSON should be two keys, status and error
	assertIntsEqual(t, "Number of keys in top-level JSON", 2, len(responseAsMap))

	// check status
	assertNotNil(t, "status field should be present", responseAsMap["status"])
	statusAsString := responseAsMap["status"].(string)
	assertStringsEqual(t, "Status field", "error", statusAsString)

	// check error
	assertNotNil(t, "error field should be present", responseAsMap["error"])
	errorAsString := responseAsMap["error"].(string)
	if expectedError != "" {
		assertStringsEqual(t, "error field", expectedError, errorAsString)
	}
}

//func TestGlobalStats(t *testing.T) {
//	// set neither start time or end time
//	//url := "http://localhost:8080/v1/stats"
//
//	// start time only - pick a time an hour ago (endtime will default to now)
//	//starttimeString := time.Now().Add(-(time.Duration(60) * time.Minute)).Format(prometheusTimeFormat)
//	//url := "http://localhost:8080/v1/stats?starttime=" + starttimeString
//
//	// end time only - pick a time an hour ago (start time will default to 5 mins before that)
//	//endtimeString := time.Now().Add(-(time.Duration(60) * time.Minute)).Format(prometheusTimeFormat)
//	//url := "http://localhost:8080/v1/stats?endtime=" + endtimeString
//
//	// start time and end time - end time is now, start time is 10 mins ago, step is 30s (so expect about 20 values)
//	endtimeString := time.Now().Format(prometheusTimeFormat)
//	starttimeString := time.Now().Add(-(time.Duration(10) * time.Minute)).Format(prometheusTimeFormat)
//	url := "http://localhost:8080/v1/stats?starttime=" + starttimeString + "&endtime=" + endtimeString + "&step=30s"
//
//	// could also test not setting step
//	// could also test invalid values for start/end/step
//	// could also test end time before start time
//
//	t.Log("Sending " + url)
//	send(t, url)
//
//}
//
//func send(t *testing.T, url string) {
//	httpClient := http.Client{
//		Timeout: time.Second * 2, // Maximum of 2 secs
//	}
//
//	req, err := http.NewRequest(http.MethodGet, url, nil)
//	if err != nil {
//		t.Fatal("Unable to perform HTTP request" + err.Error())
//		return
//	}
//
//	req.Header.Set("User-Agent", "github.com/fnproject/ext-statsapi/main-test")
//
//	res, getErr := httpClient.Do(req)
//	if getErr != nil {
//		t.Fatal("Unable to perform HTTP request" + getErr.Error())
//		return
//	}
//
//	body, readErr := ioutil.ReadAll(res.Body)
//	if readErr != nil {
//		t.Fatal("Unable to read data returned from HTTP request" + readErr.Error())
//		return
//	}
//
//	t.Log("=Received============================================")
//	t.Log(string(body[:]))
//	t.Log("=====================================================")
//
//	thisMetricsResponse := metricsResponse{}
//	jsonErr := json.Unmarshal(body, &thisMetricsResponse)
//	if jsonErr != nil {
//		t.Fatal("Error unmarshalling returned data" + jsonErr.Error())
//	}
//	t.Log("Status=" + thisMetricsResponse.Status)
//	if thisMetricsResponse.Status != "success" {
//		thisErrorResponse := errorResponse{}
//		jsonErr := json.Unmarshal(body, &thisErrorResponse)
//		if jsonErr != nil {
//			t.Fatal("Status=" + thisMetricsResponse.Status)
//		}
//		t.Fatal("Status=" + thisErrorResponse.Status + ":" + thisErrorResponse.Error)
//	}
//	t.Fatal("Failing test to force log output to be visible")
//}
