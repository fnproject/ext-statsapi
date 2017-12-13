package statistics

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

// Requires the folowing to be running
// - an extended Fn server (such as the one in this project)
// - a Prometheus server, configured to scrape data from the Fn server

func TestGlobalStats(t *testing.T) {
	// set neither start time or end time
	//url := "http://localhost:8080/v1/statistics"

	// start time only - pick a time an hour ago (endtime will default to now)
	//starttimeString := time.Now().Add(-(time.Duration(60) * time.Minute)).Format(prometheusTimeFormat)
	//url := "http://localhost:8080/v1/statistics?starttime=" + starttimeString

	// end time only - pick a time an hour ago (start time will default to 5 mins before that)
	//endtimeString := time.Now().Add(-(time.Duration(60) * time.Minute)).Format(prometheusTimeFormat)
	//url := "http://localhost:8080/v1/statistics?endtime=" + endtimeString

	// start time and end time - end time is now, start time is 10 mins ago, step is 30s (so expect about 20 values)
	endtimeString := time.Now().Format(prometheusTimeFormat)
	starttimeString := time.Now().Add(-(time.Duration(10) * time.Minute)).Format(prometheusTimeFormat)
	url := "http://localhost:8080/v1/statistics?starttime=" + starttimeString + "&endtime=" + endtimeString + "&step=30s"

	// could also test not setting step
	// could also test invalid values for start/end/step
	// could also test end time before start time

	t.Log("Sending "+url)
	send(t, url)

}

func send(t *testing.T, url string) {
	httpClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal("Unable to perform HTTP request" + err.Error())
		return
	}

	req.Header.Set("User-Agent", "github.com/fnproject/ext-metrics/main-test")

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		t.Fatal("Unable to perform HTTP request" + getErr.Error())
		return
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		t.Fatal("Unable to read data returned from HTTP request" + readErr.Error())
		return
	}


	t.Log("=Received============================================")
	t.Log(string(body[:]))
	t.Log("=====================================================")

	thisMetricsResponse := metricsResponse{}
	jsonErr := json.Unmarshal(body, &thisMetricsResponse)
	if jsonErr != nil {
		t.Fatal("Error unmarshalling returned data" + jsonErr.Error())
	}
	t.Log("Status="+thisMetricsResponse.Status)
	if thisMetricsResponse.Status != "success" {
		thisErrorResponse := errorResponse{}
		jsonErr := json.Unmarshal(body, &thisErrorResponse)
		if jsonErr != nil {
		t.Fatal("Status="+thisMetricsResponse.Status)
		}
		t.Fatal("Status="+thisErrorResponse.Status+":"+thisErrorResponse.Error)
	}
	//t.Fatal("Failing test to force log output to be visible")
}
