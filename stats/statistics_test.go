package stats

import (
	"os"
	"strconv"
	"testing"
	"time"
)

// Requires the following to be running
// - an extended Fn server (such as the one in this project)
// - a Prometheus server, configured to scrape data from the Fn server
// Also requires the following to be run beforehand
//   bash test/create.bash

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	// call some functions

	//cold sync
	for i := 0; i < 10; i++ {
		getURLWithPanic("localhost:8080/r/hello-cold-sync-a/hello-cold-sync-a1")
	}
	// hot sync
	for i := 0; i < 10; i++ {
		getURLWithPanic("http://localhost:8080/r/hello-hot-sync-a/hello-hot-sync-a1")
	}

	//cold async
	getURLWithPanic("http://localhost:8080/r/hello-cold-async-a/hello-cold-async-a1")
	getURLWithPanic("http://localhost:8080/r/hello-cold-async-a/hello-cold-async-a2")
	getURLWithPanic("http://localhost:8080/r/hello-cold-async-a/hello-cold-async-a3")
	getURLWithPanic("http://localhost:8080/r/hello-cold-async-a/hello-cold-async-a1")
	getURLWithPanic("http://localhost:8080/r/hello-cold-async-a/hello-cold-async-a2")
	getURLWithPanic("http://localhost:8080/r/hello-cold-async-a/hello-cold-async-a1")
	getURLWithPanic("http://localhost:8080/r/hello-cold-async-b/hello-cold-async-b1")
	getURLWithPanic("http://localhost:8080/r/hello-cold-async-b/hello-cold-async-b2")
	// don't call this last one so we can check the stats of a function that has never been called
	//getURLWithPanic("localhost:8080/r/hello-cold-async-b/hello-cold-async-b3")

	//hot async
	for i := 0; i < 10; i++ {
		getURLWithPanic("localhost:8080/r/hello-hot-async-a/hello-hot-async-a1")
	}
}

func shutdown() {

}

// Test a query which will return an error from the extension (rather than from the main router, such as app or route not found)
func TestBadRoute(t *testing.T) {
	url := "http://localhost:8080/v1/apps/hello-cold-sync-a/routes/hello-cold-sync-a1/stats?step=Wombat"
	response := getURLAsJSON(t, url)
	verifyFailedJSON(t, response, "Unable to parse step parameter: time: invalid duration Wombat")
}

// Test sync cold
func TestNeverCalled(t *testing.T) {
	// verify stats for a function that exists but has never been called (since server startup)
	appname := "hello-cold-async-b"
	routename := "hello-cold-async-b3"
	verifyWithRetries(t, appname, routename)
}

func TestAllFuncs(t *testing.T) {
	// verify stats for all functions
	appname := ""
	routename := ""
	verifyWithRetries(t, appname, routename)
}

func TestAllFuncsPerApp(t *testing.T) {
	// verify stats across all three functions in the app hello-cold-async-a
	appname := "hello-cold-async-a"
	routename := ""
	verifyWithRetries(t, appname, routename)
}

// Test sync cold
func TestSyncCold(t *testing.T) {
	// verify stats for sync cold functions
	appname := "hello-cold-sync-a"
	routename := "hello-cold-sync-a1"
	verifyWithRetries(t, appname, routename)
}

// Test async cold
func TestAsyncCold(t *testing.T) {
	// verify stats for async cold functions
	appname := "hello-cold-async-a"
	routename := "hello-cold-async-a1"
	verifyWithRetries(t, appname, routename)
}

// Test sync hot
func TestSyncHot(t *testing.T) {
	// verify stats for sync hot functions
	appname := "hello-hot-sync-a"
	routename := "hello-hot-sync-a1"
	verifyWithRetries(t, appname, routename)
}

// Test async hot
// Assumes test/run-hot-async.bash has been run
func TestAsyncHot(t *testing.T) {
	// verify stats for async hot functions
	appname := "hello-hot-async-a"
	routename := "hello-hot-async-a1"
	verifyWithRetries(t, appname, routename)
}

func verifyWithRetries(t *testing.T, appname string, routename string) {

	// work out what stats API URL to call
	var url string
	if appname == "" && routename == "" {
		url = "http://localhost:8080/v1/stats"
	} else if routename == "" {
		url = "http://localhost:8080/v1/apps/" + appname + "/stats"
	} else {
		url = "http://localhost:8080/v1/apps/" + appname + "/routes/" + routename + "/stats"
	}

	startTime := time.Now()
	timeout := time.Duration(60) * time.Second
	var attempt int
	var err error

	// loop until timeout
	for time.Now().Sub(startTime) < timeout {
		if attempt > 0 {
			// except for the first time, wait 10 seconds before trying again
			t.Log("Sleeping before iteration " + strconv.Itoa(attempt))
			time.Sleep(time.Duration(10) * time.Second)
		}
		attempt++

		// call the stats API and obtain the response
		response := getURLAsJSON(t, url)
		t.Log("Response on attempt " + strconv.Itoa(attempt))
		t.Log(response)

		// get basic metrics by scraping the /metrics endpoint directly
		var expectedMetrics map[string]int
		if appname == "" && routename == "" {
			expectedMetrics = getAllMetricsForAll(t)
		} else if routename == "" {
			expectedMetrics = getAllMetricsForApp(t, appname)
		} else {
			expectedMetrics = getAllMetricsForAppAndRoute(t, appname, routename)
		}

		t.Log("expectedMetrics on attempt " + strconv.Itoa(attempt))
		t.Log(expectedMetrics)

		err = verifyOnce(t, response, expectedMetrics)

		if err == nil {
			// no error encountered so break out of the loop
			break
		}

	}

	if err != nil {
		var errorMessage = "Test failed on attempt " + strconv.Itoa(attempt) + ": " + err.Error()
		t.Fatal(errorMessage)
	}

}

func verifyOnce(t *testing.T, response interface{}, expectedMetrics map[string]int) error {
	var err error

	responseAsMap := response.(map[string]interface{})

	// top level of JSON should be two keys, status and data
	err = checkIntsEqual(t, "Number of keys in top-level JSON", 2, len(responseAsMap))
	if err != nil {
		return err
	}

	// check status
	err = checkNotNil(t, "status field should be present", responseAsMap["status"])
	if err != nil {
		return err
	}
	statusAsString := responseAsMap["status"].(string)
	err = checkStringsEqual(t, "Status field", "success", statusAsString)
	if err != nil {
		return err
	}

	// check data
	err = checkNotNil(t, "data field should be present", responseAsMap["data"])
	if err != nil {
		return err
	}
	dataAsMap := responseAsMap["data"].(map[string]interface{})
	err = checkIntsEqual(t, "Number of keys in data", 6, len(dataAsMap))
	if err != nil {
		return err
	}

	// check the fields in the data array that correspond to basic metrics
	err = checkMetricField(t, dataAsMap, jsonKeys[callsConst], expectedMetrics[promMetricNames[callsConst]])
	if err != nil {
		return err
	}
	err = checkMetricField(t, dataAsMap, jsonKeys[completedConst], expectedMetrics[promMetricNames[completedConst]])
	if err != nil {
		return err
	}
	err = checkMetricField(t, dataAsMap, jsonKeys[failedConst], expectedMetrics[promMetricNames[failedConst]])
	if err != nil {
		return err
	}
	err = checkMetricField(t, dataAsMap, jsonKeys[timedoutConst], expectedMetrics[promMetricNames[timedoutConst]])
	if err != nil {
		return err
	}
	err = checkMetricField(t, dataAsMap, jsonKeys[errorsConst], expectedMetrics[promMetricNames[errorsConst]])
	if err != nil {
		return err
	}

	// check the fields in the data array that correspond to durations
	err = checkNotNil(t, "durations field should be present", dataAsMap["durations"])
	if err != nil {
		return err
	}
	durationsAsArray := dataAsMap["durations"].([]interface{})
	// if any time-value pairs are present, verify them
	for _, timeAndValuePair := range durationsAsArray {
		timeAndValuePairAsMap := timeAndValuePair.(map[string]interface{})

		// check data -> durations -> element -> time
		err = checkNotNil(t, "time field should be present", timeAndValuePairAsMap["time"])
		if err != nil {
			return err
		}

		// check data -> durations -> element -> value
		err = checkNotNil(t, "value field should be present", timeAndValuePairAsMap["value"])
		if err != nil {
			return err
		}
	}

	// success!
	return nil
}

func checkMetricField(t *testing.T, dataAsMap map[string]interface{}, jsonkey string, expectedValue int) error {
	var err error

	err = checkNotNil(t, "field "+jsonkey+" should be present", dataAsMap[jsonkey])
	if err != nil {
		return err
	}
	metricAsArray := dataAsMap[jsonkey].([]interface{})
	// if any time-value pairs are present, verify them
	for _, timeAndValuePair := range metricAsArray {
		timeAndValuePairAsMap := timeAndValuePair.(map[string]interface{})
		// check data -> completed -> time
		err = checkNotNil(t, "time field should be present", timeAndValuePairAsMap["time"])
		if err != nil {
			return err
		}
		// check data -> completed -> value
		err = checkNotNil(t, "value field should be present", timeAndValuePairAsMap["value"])
		if err != nil {
			return err
		}
	}
	// verify latest metric count matches expectedValue
	if len(metricAsArray) > 0 {
		lastTimeAndValuePair := metricAsArray[len(metricAsArray)-1]
		lastTimeAndValuePairAsMap := lastTimeAndValuePair.(map[string]interface{})
		lastValue := lastTimeAndValuePairAsMap["value"]
		lastValueAsFloat64 := lastValue.(float64)
		lastvalueAsInt := int(lastValueAsFloat64)
		err = checkIntsEqual(t, jsonkey, expectedValue, lastvalueAsInt)
		if err != nil {
			return err
		}
	} else {
		// metric array is empty which implies a ount of zero
		err = checkIntsEqual(t, jsonkey+" array is empty", expectedValue, 0)
		if err != nil {
			return err
		}
	}
	return nil
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
