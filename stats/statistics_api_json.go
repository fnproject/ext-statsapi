package statistics

// the following structs represent the JSON that we are building and returning to our caller
// we have different structs for error and success

const STATS_STATUS_ERROR string = "error"
const STATS_STATUS_SUCCESS string = "success"

type errorResponse struct {
	Status string `json:"status"` // "error" (STATS_STATUS_ERROR)
	Error  string `json:"error"`  //  if Status is "error", set to the error message
}

type metricsResponse struct {
	Status string                            `json:"status"` // "success" (STATS_STATUS_SUCCESS)
	Data   map[string][]metricsTimeValuePair `json:"data"`
}

type metricsTimeValuePair struct {
	Time  int64   `json:"time"`
	Value float64 `json:"value"`
}
