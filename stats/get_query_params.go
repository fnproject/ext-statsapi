package statistics

import (
	"errors"
	"net/http"
	"time"
)

// time format required when sending queries to Prometheus
const prometheusTimeFormat = "2006-01-02T15:04:05.999Z07:00"

// Extract and return the required URL query parameters, generating default values if missing
func getQueryParams(r *http.Request) (string, string, string, error) {

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
	case len(startTimeParams) == 0:
		// endtime is specified, starttime is not specified, set to 5mins before endtime
		starttime = endtime.Add(-(time.Duration(5) * time.Minute))
		starttimeString = starttime.Format(prometheusTimeFormat)
	case len(endTimeParams) == 0:
		// starttime is specified, endtime is not specified, set to now
		endtime = time.Now()
		endtimeString = endtime.Format(prometheusTimeFormat)
	default:
		// both starttime and endtime specified
	}

	if endtime.Before(starttime) {
		err = errors.New("endtime (" + endtimeString + ") is before starttime (" + starttimeString + ")")
		return "", "", "", err
	}

	if len(stepParams) == 0 {
		// step not specified, assume 30 secs
		stepString = time.Duration(30 * time.Second).String()
	}

	return starttimeString, endtimeString, stepString, err
}