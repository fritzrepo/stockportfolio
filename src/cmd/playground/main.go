package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type ApiResponse struct {
	Status       string `json:"status"`
	Message      string `json:"message"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	ErrorDetails string `json:"errorDetails,omitempty"`
	Data         any    `json:"data,omitempty"`
}

func main() {
	log.Println("Playgroud initialized successfully.")

	// var tmp *ApiResponse

	// tmp = &ApiResponse{
	// 	Message: "Hallo",
	// }

	// tmp2 := tmp

	// log.Println("Temporary response:", tmp.Message)
	// tmp.Message = "Hallo2"
	// log.Println("Temporary response:", tmp.Message)
	// log.Println("Temporary response2:", tmp2.Message)

	sessionId := uuid.New().String()

	var requestId string
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	timestampStr := fmt.Sprintf("%d", timestamp)
	requestId = timestampStr[len(timestampStr)-9:]

	fmt.Println("Session ID:", sessionId)
	fmt.Println("Request ID:", requestId)

	infoValue := fmt.Sprintf("{\"clientRequestId\":{\"sessionId\":\"%s\",\"requestId\":\"%s\"}}", sessionId, requestId)
	log.Println("Info Value:", infoValue)

}
