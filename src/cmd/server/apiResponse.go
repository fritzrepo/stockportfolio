package main

type ApiResponse struct {
	Status       string `json:"status"`
	Message      string `json:"message"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	ErrorDetails string `json:"errorDetails,omitempty"`
	Data         any    `json:"data,omitempty"`
}
