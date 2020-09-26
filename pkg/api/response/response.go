package response

// HTTPErr represents an error response for the server
type HTTPErr struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

// StatusResponse represents the schema of the /internal/status response
type StatusResponse struct {
	Running bool `json:"running"`
}
