package stats

// Prometheus metrics to use, keyed by metric type
// see comment in statistics.go for information on adding a new metric
var promMetricNames = map[int]string{
	completedConst: "fn_completed",
	failedConst:    "fn_failed",
	callsConst:     "fn_calls",
	errorsConst:    "fn_errors",
	timedoutConst:  "fn_timeouts",
	runningConst:   "fn_running", // used for tests only
	queuedConst:    "fn_queued",  // used for tests only
	durationsConst: "fn_span_agent_submit_duration_seconds",
}

// Functions that know how to build the required Prometheus query, keyed by metric type
// see comment in statistics.go for information on adding a new metric
var queryBuilders = map[int]func(string, string, string, string, string, string, string, string) string{
	completedConst: queryBuilderForCountersAndGauges,
	failedConst:    queryBuilderForCountersAndGauges,
	callsConst:     queryBuilderForCountersAndGauges,
	errorsConst:    queryBuilderForCountersAndGauges,
	timedoutConst:  queryBuilderForCountersAndGauges,
	durationsConst: queryBuilderForHistograms,
}

func buildPrometheusRequest(queryBuilder func(string, string, string, string, string, string, string, string) string, promHost string, promPort string, metricType int, appName string, routeName string, startTimeString string, endTimeString string, stepString string) string {
	promMetricName := promMetricNames[metricType]
	// use the specified queryBuilder function to construct a Prometheus query for the required metric
	query := queryBuilder(promHost, promPort, promMetricName, appName, routeName, startTimeString, endTimeString, stepString)
	// construct the complete request URL, including host, port, time range and step
	return "http://" + promHost + ":" + promPort + "/api/v1/query_range?query=" + query + "&start=" + startTimeString + "&end=" + endTimeString + "&step=" + stepString
}

func queryBuilderForCountersAndGauges(promHost string, promPort string, promMetricName string, appName string, routeName string, startTimeString string, endTimeString string, stepString string) string {

	if appName == "" {
		return "sum(" + promMetricName + ")"
	} else if routeName == "" {
		qualifiedPromMetricName := promMetricName + "{" + appLabel + "=\"" + appName + "\"}"
		return "sum(" + qualifiedPromMetricName + ")"
	} else {
		qualifiedPromMetricName := promMetricName + "{" + appLabel + "=\"" + appName + "\"," + routeLabel + "=\"" + routeName + "\"}"
		return "sum(" + qualifiedPromMetricName + ")"
	}

}

func queryBuilderForHistograms(promHost string, promPort string, promMetricName string, appName string, routeName string, startTimeString string, endTimeString string, stepString string) string {

	rollingMeanPeriod := "1m"
	if appName == "" {
		numerator := "sum(rate(" + promMetricName + "_sum[" + rollingMeanPeriod + "]))"
		denominator := "sum(rate(" + promMetricName + "_count[" + rollingMeanPeriod + "]))"
		return numerator + "/" + denominator
	} else if routeName == "" {
		numerator := "sum(rate(" + promMetricName + "_sum{" + appLabel + "=\"" + appName + "\"}[" + rollingMeanPeriod + "]))"
		denominator := "sum(rate(" + promMetricName + "_count{" + appLabel + "=\"" + appName + "\"}[" + rollingMeanPeriod + "]))"
		return numerator + "/" + denominator
	} else {
		numerator := "sum(rate(" + promMetricName + "_sum{" + appLabel + "=\"" + appName + "\"," + routeLabel + "=\"" + routeName + "\"}[" + rollingMeanPeriod + "]))"
		denominator := "sum(rate(" + promMetricName + "_count{" + appLabel + "=\"" + appName + "\"," + routeLabel + "=\"" + routeName + "\"}[" + rollingMeanPeriod + "]))"
		return numerator + "/" + denominator
	}
}
