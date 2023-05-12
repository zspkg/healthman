package resources

// ServiceHealthResponse is a response for /health endpoint containing outcome
type ServiceHealthResponse struct {
	Outcome State   `json:"outcome"`
	Checks  []Check `json:"checks"`
}

// Check is a structure representing health of a service
type Check struct {
	Name  string `json:"name"`
	State State  `json:"state"`
}
