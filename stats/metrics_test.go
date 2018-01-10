package stats

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

// name of Prometheus metrics
var callsMet = "fn_api_calls"
var queuedMet = "fn_api_queued"
var runningMet = "fn_api_running"
var failedMet = "fn_api_failed"
var completedMet = "fn_api_completed"

// Test Prometheus metrics directly

// Requires the following to be running
// - either a normal or an extended Fn server (such as the one in this project)
// - no need for Prometheus to be running (these tests scrape the Fn server's /metrics endpoint directly)
// Also requires the following to be run beforehand
//   bash test/create.bash

func TestColdSyncSuccessful(t *testing.T) {

	appname := "hello-cold-sync-a"
	routename := "hello-cold-sync-a1"
	sync := true
	doTestSuccessful(t, appname, routename, sync)
}

func TestColdSyncWithTimeout(t *testing.T) {

	appname := "hello-cold-sync-a"
	routename := "hello-cold-sync-a1"
	sync := true
	hot := false
	doTestWithTimeout(t, appname, routename, sync, hot)
}

func TestColdSyncWithPanic(t *testing.T) {

	appname := "hello-cold-sync-a"
	routename := "hello-cold-sync-a1"
	sync := true
	doTestWithPanic(t, appname, routename, sync)
}

func TestColdAsyncSuccessful(t *testing.T) {

	appname := "hello-cold-async-a"
	routename := "hello-cold-async-a1"
	sync := false
	doTestSuccessful(t, appname, routename, sync)

}

func TestColdAsyncWithTimeout(t *testing.T) {

	appname := "hello-cold-async-a"
	routename := "hello-cold-async-a1"
	sync := false
	hot := false
	doTestWithTimeout(t, appname, routename, sync, hot)
}

func TestColdAsyncWithPanic(t *testing.T) {

	appname := "hello-cold-async-a"
	routename := "hello-cold-async-a1"
	sync := false
	doTestWithPanic(t, appname, routename, sync)
}

func TestHotSyncSuccessful(t *testing.T) {

	appname := "hello-hot-sync-a"
	routename := "hello-hot-sync-a1"
	sync := true
	doTestSuccessful(t, appname, routename, sync)
}

func TestHotSyncWithTimeout(t *testing.T) {

	appname := "hello-hot-sync-a"
	routename := "hello-hot-sync-a1"
	sync := true
	hot := true
	doTestWithTimeout(t, appname, routename, sync, hot)
}

func TestHotSyncWithPanic(t *testing.T) {

	appname := "hello-hot-sync-a"
	routename := "hello-hot-sync-a1"
	sync := true
	doTestWithPanic(t, appname, routename, sync)
}

func TestHotAsyncSuccessful(t *testing.T) {

	appname := "hello-hot-async-a"
	routename := "hello-hot-async-a1"
	sync := false
	doTestSuccessful(t, appname, routename, sync)
}

func TestHotAsyncWithTimeout(t *testing.T) {

	appname := "hello-hot-async-a"
	routename := "hello-hot-async-a1"
	sync := false
	hot := true
	doTestWithTimeout(t, appname, routename, sync, hot)
}

func TestHotAsyncWithPanic(t *testing.T) {

	appname := "hello-hot-async-a"
	routename := "hello-hot-async-a1"
	sync := false
	doTestWithPanic(t, appname, routename, sync)
}

func doTestSuccessful(t *testing.T, appname string, routename string, sync bool) {

	metrics0 := getMetrics(t, appname, routename)

	// make a function call which will NOT timeout
	sleeptime := 0
	forceTimeout := false
	forcePanic := false
	output := call(t, appname, routename, sync, forceTimeout, forcePanic)
	if !strings.Contains(output, "COMPLETEDOK") {
		t.Fatal("Function did not complete: " + output)
	}

	metrics1 := getMetrics(t, appname, routename)

	message := "after calling " + appname + "/" + routename + " with sleeptime " + strconv.Itoa(sleeptime)
	assertIntsEqual(t, message+" calls should have increased by 1", metrics0[callsMet]+1, metrics1[callsMet])
	assertIntsEqual(t, message+" completed should have incrreturnseased by 1", metrics0[completedMet]+1, metrics1[completedMet])
	assertIntsEqual(t, message+" queued should be unchanged", metrics0[queuedMet], metrics1[queuedMet])
	assertIntsEqual(t, message+" failed should be unchanged", metrics0[failedMet], metrics1[failedMet])
	assertIntsEqual(t, message+" running should be unchanged", metrics0[runningMet], metrics1[runningMet])

}

func doTestWithTimeout(t *testing.T, appname string, routename string, sync bool, hot bool) {

	metrics0 := getMetrics(t, appname, routename)

	// make a function call which WILL timeout"
	sleeptime := 60000 // Function timeout is 5000 ms
	forceTimeout := true
	forcePanic := false
	output := call(t, appname, routename, sync, forceTimeout, forcePanic)
	if sync {
		if !strings.Contains(output, "Timed out") {
			t.Fatal("Function call did not return system-generated timeout message as expected: " + output)
		}
		if hot {
			// not sure what is supposed to happen, skip testing
		} else {
			if !strings.Contains(output, "FORCETIMEOUT") {
				t.Fatal("Function call did not return function output: " + output)
			}
		}
	} else {
		if hot {
			// now sure what is supposed to happen, skip testing
		} else {
			if strings.Contains(output, "Timed out") {
				t.Fatal("Function call unexpectedly returned system-generated timeout message: " + output)
			}
			if !strings.Contains(output, "FORCETIMEOUT") {
				t.Fatal("Function call does not return function output: " + output)
			}
		}
	}

	metrics1 := getMetrics(t, appname, routename)

	message := "after calling " + appname + "/" + routename + " with sleeptime " + strconv.Itoa(sleeptime)
	assertIntsEqual(t, message+" calls should have increased by 1", metrics0[callsMet]+1, metrics1[callsMet])
	assertIntsEqual(t, message+" completed should be unchanged", metrics0[completedMet], metrics1[completedMet])
	assertIntsEqual(t, message+" queued should be unchanged", metrics0[queuedMet], metrics1[queuedMet])
	assertIntsEqual(t, message+" failed should have increased by 1", metrics0[failedMet]+1, metrics1[failedMet])
	assertIntsEqual(t, message+" running should be unchanged", metrics0[runningMet], metrics1[runningMet])

}

func doTestWithPanic(t *testing.T, appname string, routename string, sync bool) {

	metrics0 := getMetrics(t, appname, routename)

	// make a function call which will panic"
	sleeptime := 0
	forceTimeout := false
	forcePanic := true
	output := call(t, appname, routename, sync, forceTimeout, forcePanic)
	if sync {
		// we find by experiment that the function panic output is lost, but the system returns some JSON containing the message "container exit code 2"
		// This is just how things happen to be
		if !strings.Contains(output, "container exit code 2") {
			t.Fatal("Function call did not return system-generated timeout message as expected: " + output)
		}
		if strings.Contains(output, "FORCEPANIC") {
			t.Fatal("Function call unexpectedly returned function output: " + output)
		}
	} else {
		// we find by experiment that the function panic output is available, and the system does not return the message "container exit code 2"
		// This is just how things happen to be
		if strings.Contains(output, "container exit code 2") {
			t.Fatal("Function call unexpectedly returned system-generated timeout message: " + output)
		}
		if !strings.Contains(output, "panic: FORCEPANIC") {
			t.Fatal("Function call does not return function output: " + output)
		}
	}

	metrics1 := getMetrics(t, appname, routename)

	message := "after calling " + appname + "/" + routename + " with sleeptime " + strconv.Itoa(sleeptime)
	assertIntsEqual(t, message+" calls should have increased by 1", metrics0[callsMet]+1, metrics1[callsMet])
	assertIntsEqual(t, message+" completed should be unchanged", metrics0[completedMet], metrics1[completedMet])
	assertIntsEqual(t, message+" queued should be unchanged", metrics0[queuedMet], metrics1[queuedMet])
	assertIntsEqual(t, message+" failed should have increased by 1", metrics0[failedMet]+1, metrics1[failedMet])
	assertIntsEqual(t, message+" running should be unchanged", metrics0[runningMet], metrics1[runningMet])

}

func call(t *testing.T, appname string, routename string, sync bool, forceTimeout bool, forcePanic bool) string {

	// if both forceTimeout and forcePanic are set then function should panic

	url := "http://localhost:8080/r/" + appname + "/" + routename +
		"?forcetimeout=" + strconv.FormatBool(forceTimeout) + "&forcepanic=" + strconv.FormatBool(forcePanic)

	var response string
	if sync {
		response = getURLAsString(t, url)
		if response == `{"error":{"message":"Timed out"}}"` {
			t.Fatal(errors.New("Sync call timed out"))
		}
	} else {
		jsonResponse := getURLAsJSON(t, url)
		responseAsMap := jsonResponse.(map[string]interface{})
		callid := responseAsMap["call_id"].(string)
		// now wait for the async call to complete
		statusQueryURL := "http://localhost:8080/v1/apps/" + appname + "/calls/" + callid
		completed := false
		startTime := time.Now()
		timeout := 60 * time.Second
		for !completed {
			// wait for the call to complete
			callStatusJSON := getURLAsJSON(t, statusQueryURL)
			callStatusAsMap := callStatusJSON.(map[string]interface{})
			callJSON := callStatusAsMap["call"]
			if callJSON != nil {
				callAsMap := callJSON.(map[string]interface{})
				println(callAsMap)
				status := callAsMap["status"]
				statusAsString := status.(string)
				var expectedStatus string
				if forceTimeout && !forcePanic {
					expectedStatus = "timeout"
				} else if forcePanic {
					expectedStatus = "error"
				} else {
					expectedStatus = "success"
				}
				if statusAsString == expectedStatus {
					completed = true
				}
			}
			if !completed {
				if time.Now().Sub(startTime) > timeout {
					t.Fatal(errors.New("Async call appears to have not been executed"))
				}
				time.Sleep(time.Duration(1000) * time.Millisecond)
			}
		}
		// async call has completed, now get the output
		jsonLogResponse := getURLAsJSON(t, statusQueryURL+"/log")
		logResponseAsMap := jsonLogResponse.(map[string]interface{})
		logAsJSON := logResponseAsMap["log"]
		logAsMap := logAsJSON.(map[string]interface{})
		loglogAsJSON := logAsMap["log"]
		response = loglogAsJSON.(string)

	}

	// now check the status of the call

	return response
}

func getMetrics(t *testing.T, appname string, routename string) map[string]int {

	result := make(map[string]int)

	// get all Prometheus metrics
	scrapedMetrics := getURLAsString(t, "http://localhost:8080/metrics")

	requiredMetrics := []string{"fn_api_calls", "fn_api_queued", "fn_api_completed", "fn_api_failed", "fn_api_running"}
	for _, thisMetricName := range requiredMetrics {
		var thisMetricValue int
		var err error
		regularExpression := thisMetricName + `{app="` + appname + `",path="/` + routename + `"} (\d+)`
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
