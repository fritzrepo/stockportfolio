package main

import "log"

type ApiResponse struct {
	Status       string `json:"status"`
	Message      string `json:"message"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	ErrorDetails string `json:"errorDetails,omitempty"`
	Data         any    `json:"data,omitempty"`
}

func main() {
	log.Println("Playgroud initialized successfully.")

	var tmp *ApiResponse

	tmp = &ApiResponse{
		Message: "Hallo",
	}

	tmp2 := tmp

	log.Println("Temporary response:", tmp.Message)
	tmp.Message = "Hallo2"
	log.Println("Temporary response:", tmp.Message)
	log.Println("Temporary response2:", tmp2.Message)

}
