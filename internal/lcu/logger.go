package lcu

// APILogEntry represents a single LCU API call for logging/telemetry.
// Duration is expressed in milliseconds for easy display in the frontend.
type APILogEntry struct {
	Type       string            `json:"type"`               // always "lcu"
	Method     string            `json:"method"`             // GET, POST, etc.
	Endpoint   string            `json:"endpoint"`           // e.g. /lol-matchmaking/v1/ready-check
	StatusCode int               `json:"statusCode"`         // HTTP status code
	Duration   int64             `json:"duration"`           // request duration in ms
	Headers    map[string]string `json:"headers,omitempty"`  // request headers
	Response   string            `json:"response,omitempty"` // optional response body (truncated)
	Error      string            `json:"error,omitempty"`
}

// apiLogger is an optional callback set by the app to receive API logs.
var apiLogger func(entry APILogEntry)

// SetAPILogger registers a callback to receive LCU API logs.
func SetAPILogger(logger func(entry APILogEntry)) {
	apiLogger = logger
}

// logAPICall sends an API log entry to the registered logger if present.
func logAPICall(entry APILogEntry) {
	if apiLogger != nil {
		apiLogger(entry)
	}
}
