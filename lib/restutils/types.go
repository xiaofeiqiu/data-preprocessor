package restutils

type HealthResponse struct {
	Status string `json:status`
}

type ErrorResponse struct {
	Event string `json:event`
	Error string `json:error`
}

func NewErrorResponse(event string, body []byte) ErrorResponse {
	return ErrorResponse{
		Event: event,
		Error: string(body),
	}
}
