package main

import (
	"fmt"
	"log"
	"time"

	"github.com/fritzrepo/stockportfolio/internal/banking/comdirect"
)

func main() {

	comm, err := comdirect.GetCommunication("./comdirectConfig.json", "./.comdirectCredentials.json")
	if err != nil {
		log.Println("Error getting communication service:", err)
		return
	}

	log.Println("Get access token for comdirect API.")
	message, err := comm.StartSession()
	if err != nil {
		log.Println("Error starting session:", err)
		return
	}
	defer comm.EndSession()

	fmt.Println("Session started successfully.")
	fmt.Println("Message:", message)

	fmt.Println("Starting periodic refresh of access token every 9 minutes...")
	comm.RefreshTokenPeriodically(9 * time.Minute)

	fmt.Println("\nPress Enter to exit...")
	fmt.Scanln()
}
