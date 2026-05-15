package responses

// GetHealthResponse is the gateway liveness payload.
type GetHealthResponse struct {
	Status string `json:"status"`
}
