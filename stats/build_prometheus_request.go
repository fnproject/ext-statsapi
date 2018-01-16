package stats

import (
	"strconv"
)

// Prometheus metrics to use for global-scope queries, keyed by metric type
// see comment in statistics.go for information on adding a new metric
var promMetricNamesForGlobalQueries = map[int]string{
	completedConst: "fn_completed",
	failedConst:    "fn_failed",
	callsConst:     "fn_calls",
	errorsConst:    "fn_errors",
	timedoutConst:  "fn_timedout",
	durationsConst: "fn_span_agent_submit_global_duration_seconds",
}

// Prometheus metrics to use for app-scoped queries, keyed by metric type
// see comment in statistics.go for information on adding a new metric
var promMetricNamesForAppScopedQueries = map[int]string{
	completedConst: "fn_completed",
	failedConst:    "fn_failed",
	callsConst:     "fn_calls",
	errorsConst:    "fn_errors",
	timedoutConst:  "fn_timedout",
	durationsConst: "fn_span_agent_submit_app_duration_seconds",
}

// Prometheus metrics to use for route-scoped queries, keyed by metric type
// see comment in statistics.go for information on adding a new metric
var promMetricNamesForRouteScopedQueries = map[int]string{
	completedConst: "fn_completed",
	failedConst:    "fn_failed",
	callsConst:     "fn_calls",
	errorsConst:    "fn_errors",
	timedoutConst:  "fn_timedout",
	durationsConst: "fn_span_agent_submit_duration_seconds",
}

// Functions that know how to build the required Prometheus query, keyed by metric type
// see comment in statistics.go for information on adding a new metric
var queryBuilders = map[int]func(string, string, string, int, int, string, string, string, string, string) string{
	completedConst: queryBuilderForCountersAndGauges,
	failedConst:    queryBuilderForCountersAndGauges,
	callsConst:     queryBuilderForCountersAndGauges,
	errorsConst:    queryBuilderForCountersAndGauges,
	timedoutConst:  queryBuilderForCountersAndGauges,
	durationsConst: queryBuilderForForDurations,
}

var promMetricNameMapsForQueries = make(map[int]map[int]string)

func init() {
	promMetricNameMapsForQueries[query_scope_global] = promMetricNamesForGlobalQueries
	promMetricNameMapsForQueries[query_scope_app] = promMetricNamesForAppScopedQueries
	promMetricNameMapsForQueries[query_scope_route] = promMetricNamesForRouteScopedQueries
}

func buildPrometheusRequest(queryBuilder func(string, string, string, int, int, string, string, string, string, string) string, promHost string, promPort string, queryScope int, metricType int, appName string, routeName string, startTimeString string, endTimeString string, stepString string) string {
	promMetricName := promMetricNameMapsForQueries[queryScope][metricType]
	// use the specified queryBuilder function to construct a Prometheus query for the required metric
	query := queryBuilder(promHost, promPort, promMetricName, queryScope, metricType, appName, routeName, startTimeString, endTimeString, stepString)
	// construct the complete request URL, including host, port, time range and step
	return "http://" + promHost + ":" + promPort + "/api/v1/query_range?query=" + query + "&start=" + startTimeString + "&end=" + endTimeString + "&step=" + stepString
}

func queryBuilderForCountersAndGauges(promHost string, promPort string, promMetricName string, queryScope int, metricType int, appName string, routeName string, startTimeString string, endTimeString string, stepString string) string {

	var query string
	appLabel := "fn_appname"
	routeLabel := "fn_path"
	switch queryScope {
	case query_scope_global:
		query = "sum(" + promMetricName + ")"
	case query_scope_app:
		qualifiedPromMetricName := promMetricName + "{" + appLabel + "=\"" + appName + "\"}"
		query = "sum(" + qualifiedPromMetricName + ")"
	case query_scope_route:
		qualifiedPromMetricName := promMetricName + "{" + appLabel + "=\"" + appName + "\"," + routeLabel + "=\"" + routeName + "\"}"
		query = "sum(" + qualifiedPromMetricName + ")"
	default:
		panic("Unexpected queryScope" + strconv.Itoa(queryScope))
	}
	return query
}

func queryBuilderForForDurations(promHost string, promPort string, promMetricName string, queryScope int, metricType int, appName string, routeName string, startTimeString string, endTimeString string, stepString string) string {

	var query string
	appLabel := "fn_appname"
	routeLabel := "fn_path"
	rollingMeanPeriod := "1m"
	switch queryScope {
	case query_scope_global:
		numerator := "rate(" + promMetricName + "_sum[" + rollingMeanPeriod + "])"
		denominator := "rate(" + promMetricName + "_count[" + rollingMeanPeriod + "])"
		query = numerator + "/" + denominator
	case query_scope_app:
		numerator := "rate(" + promMetricName + "_sum{" + appLabel + "=\"" + appName + "\"}[" + rollingMeanPeriod + "])"
		denominator := "rate(" + promMetricName + "_count{" + appLabel + "=\"" + appName + "\"}[" + rollingMeanPeriod + "])"
		query = numerator + "/" + denominator
	case query_scope_route:
		numerator := "rate(" + promMetricName + "_sum{" + appLabel + "=\"" + appName + "\"," + routeLabel + "=\"" + routeName + "\"}[" + rollingMeanPeriod + "])"
		denominator := "rate(" +
			promMetricName + "_count{" + appLabel + "=\"" + appName + "\"," + routeLabel + "=\"" + routeName + "\"}[" + rollingMeanPeriod + "])"
		query = numerator + "/" + denominator
	default:
		panic("Unexpected queryScope" + strconv.Itoa(queryScope))
	}
	return query
}
