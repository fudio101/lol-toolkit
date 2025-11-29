package lol

// APILogEntry represents a single Riot API call for logging/telemetry.
// Duration is expressed in milliseconds for easy display in the frontend.
type APILogEntry struct {
	Type       string            `json:"type"`               // always "riot"
	Method     string            `json:"method"`             // e.g. GET
	Endpoint   string            `json:"endpoint"`           // logical endpoint name or path
	StatusCode int               `json:"statusCode"`         // synthetic: 200 on success, 500 on error
	Duration   int64             `json:"duration"`           // call duration in ms
	Headers    map[string]string `json:"headers,omitempty"`  // request headers
	Response   string            `json:"response,omitempty"` // optional response body (JSON)
	Error      string            `json:"error,omitempty"`
}

// apiLogger is an optional callback set by the app to receive API logs.
var apiLogger func(entry APILogEntry)

// SetAPILogger registers a callback to receive Riot API logs.
func SetAPILogger(logger func(entry APILogEntry)) {
	apiLogger = logger
}

// logAPICall sends an API log entry to the registered logger if present.
func logAPICall(entry APILogEntry) {
	if apiLogger != nil {
		apiLogger(entry)
	}
}
