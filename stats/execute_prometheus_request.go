package stats

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// Use the specified URL to get a range of data values for a single metric and return it as an array of time-value pairs
func executePrometheusRequest(url string) ([]metricsTimeValuePair, error) {

	promClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "github.com/fnproject/ext-statsapi")

	res, doErr := promClient.Do(req)
	if doErr != nil {
		return nil, doErr
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	// Assume result is of type "matrix" (meaning this is a range vector)

	thisPromQueryRangeData := promQueryRangeData{}
	jsonErr := json.Unmarshal(body, &thisPromQueryRangeData)
	if jsonErr != nil {
		return nil, jsonErr
	}

	if thisPromQueryRangeData.Status != "success" {
		return nil, errors.New("Error from Prometheus: " + thisPromQueryRangeData.ErrorType + ": " + thisPromQueryRangeData.Error)
	}

	if len(thisPromQueryRangeData.Data.Result) > 1 {
		// Range query has returned multiple ranges! This should never happen: we must have got the query wrong
		// Return a suitably verbose error message to allow investigation
		return nil, errors.New("data array returned by Prometheus has more than one element: url=" + url + ", returned JSON=" + string(body[:]))
	}

	// how many data time-value pairs have we been given?
	var numberOfTimeValuePairs int
	if len(thisPromQueryRangeData.Data.Result) == 0 {
		numberOfTimeValuePairs = 0
	} else {
		numberOfTimeValuePairs = len(thisPromQueryRangeData.Data.Result[0].Value)
	}

	// make the returned array for this metric
	metricDataArray := make([]metricsTimeValuePair, numberOfTimeValuePairs)

	if numberOfTimeValuePairs > 0 {
		countOfNonNanValues := 0
		// populate the array with zero or more metricsTimeValuePair strucs
		for _, val := range thisPromQueryRangeData.Data.Result[0].Value {
			tvp := new(metricsTimeValuePair)
			tvp.Time = int64(val.UnixTime())
			// filter out NaN values
			if val.ScalarValue() != "NaN" {
				value, err := strconv.ParseFloat(val.ScalarValue(), 64)
				tvp.Value = value
				if err != nil {
					return nil, errors.New("Error converting " + val.ScalarValue() + " to a float64")
				}
				metricDataArray[countOfNonNanValues] = *tvp
				countOfNonNanValues++
			}
		}
		return metricDataArray[0:countOfNonNanValues], nil
	}
	return metricDataArray, nil

}
