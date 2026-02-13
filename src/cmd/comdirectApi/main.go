package main

import (
	"fmt"
	"log"
	"time"

	"github.com/fritzrepo/stockportfolio/internal/banking/comdirect"
)

func main() {

	comm, err := comdirect.NewCommunication("./comdirectConfig.json", "./.comdirectCredentials.json")
	if err != nil {
		log.Println("Error getting communication service:", err)
		return
	}

	fmt.Println("\nPress Enter to starting session....")
	fmt.Scanln()

	err = comm.StartSession()
	if err != nil {
		log.Println("Error starting session:", err)
		return
	}
	defer comm.EndSession()

	fmt.Println("Session started successfully.")

	fmt.Println("Starting periodic refresh of access token every 2 minutes...")
	comm.RefreshTokenPeriodically(2 * time.Minute)

	fmt.Println("\nPress Enter to exit to  the app...")
	fmt.Scanln()
}
