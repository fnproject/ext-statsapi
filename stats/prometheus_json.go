package stats

// the following structs represent the JSON returned by Prometheus API from a query_range query

const PROM_STATUS_ERROR string = "error"
const PROM_STATUS_SUCCESS string = "success"

type promQueryRangeData struct {
	Status    string `json:"status"`    // "success" | "error" | ?
	Error     string `json:"errorType"` // only present if status=error
	ErrorType string `json:"error"`     // only present if status=error
	Data      data   `json:"data"`
}

type data struct {
	ResultType string         `json:"resultType"` // "matrix" (Range vectors) | "vector" (Instant vectors) | "scalar" (Scalar results) | "string" (String results)
	Result     []matrixResult `json:"result"`
}

type matrixResult struct { // Used when resultType is matrix
	Metric map[string]string `json:"metric"` // Map of label_name to label_value (empty if this is a sum)
	Value  []timeValuePair   `json:"values"`
}

type timeValuePair []interface{} // This is an array with two elements, <unix_time>, "<scalar_value>"

func (tvp timeValuePair) UnixTime() float64 {
	result, _ := tvp[0].(float64)
	return result
}

func (tvp timeValuePair) ScalarValue() string {
	result, _ := tvp[1].(string)
	return result
}
