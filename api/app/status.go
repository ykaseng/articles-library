package app

// Status holds specific application settings linked to a Status.
type Status struct {
	Code    int    `json:"status"` // http response status code
	Message string `json:"mesage"` // user-level status message
}
