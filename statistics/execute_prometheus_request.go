package statistics

import (
	"encoding/json"
	"errors"
	"fmt"
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

	//	fmt.Println("== URL sent to Prometheus= ===========================================")
	//	fmt.Println(url)
	//	fmt.Println("===End of URL sent to Prometheus =====================================")

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

	//	fmt.Println("== JSON returned from Prometheus= ====================================")
	//	fmt.Println(string(body[:]))
	//	fmt.Println("===End of JSON returned from Prometheus ==============================")

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
	} else {
		return metricDataArray, nil
	}
}
