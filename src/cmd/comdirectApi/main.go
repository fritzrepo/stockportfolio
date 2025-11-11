package main

import (
	"fmt"
	"log"

	"github.com/fritzrepo/stockportfolio/internal/banking/comdirect"
)

func main() {

	log.Println("Start comdirect API.")
	comm, err := comdirect.GetCommunication("./comdirectConfig.json", "./.comdirectCredentials.json")
	if err != nil {
		log.Println("Error getting communication service:", err)
		return
	}

	message, err := comm.StartSession()
	if err != nil {
		log.Println("Error starting session:", err)
		return
	}
	defer comm.EndSession()

	fmt.Println("Session started successfully.")
	fmt.Println("Message:", message)
	fmt.Println("\nPress Enter to exit...")
	fmt.Scanln()
}
