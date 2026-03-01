package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fritzrepo/stockportfolio/internal/banking/comdirect"
)

func main() {
	var balances = false
	var depot = false
	var instrument = false
	var orders = false

	// the first argument is always program name
	argLength := len(os.Args[1:])
	fmt.Printf("Arg length is %d\n", argLength)

	for i, a := range os.Args[1:] {
		fmt.Printf("Arg %d is %s\n", i+1, a)
		if a == "balances" {
			balances = true
		}
		if a == "depot" {
			depot = true
		}
		if a == "instrument" {
			instrument = true
		}
		if a == "orders" {
			orders = true
		}
	}

	comm, err := comdirect.NewCommunication("./comdirectConfig.json", "./.comdirectCredentials.json")
	if err != nil {
		log.Println("Error getting communication service:", err)
		return
	}

	fmt.Println("Press Enter to starting session....")
	fmt.Scanln()

	err = comm.StartSession()
	if err != nil {
		log.Println("Error starting session:", err)
		return
	}
	defer comm.EndSession()

	fmt.Println("Session started successfully.")

	// fmt.Println("Starting periodic refresh of access token every 2 minutes...")
	// comm.RefreshTokenPeriodically(2 * time.Minute)

	if balances {

		// Fetch all account balances
		fmt.Println("Press Enter to fetching all account balances...")
		fmt.Scanln()
		accountBalances, err := comm.GetAccountBalances()
		if err != nil {
			log.Println("Error fetching account balances:", err)
			return
		}

		for _, balance := range accountBalances {
			fmt.Printf("Account ID: %s, Balance: %s %s\n", balance.Account.AccountType.Text, balance.Balance.Value, balance.Balance.Unit)
		}

		var accountId string
		accountId = accountBalances[0].Account.AccountId

		// Fetch balance for specific account
		fmt.Printf("Press Enter to fetching balance for account ID: %s...\n", accountId)
		fmt.Scanln()
		accountBalance, err := comm.GetAccoutBalance(accountId)
		if err != nil {
			log.Println("Error fetching account balance:", err)
			return
		}
		fmt.Printf("Account ID: %s, Balance: %s %s\n", accountBalance.Account.AccountType.Text, accountBalance.Balance.Value, accountBalance.Balance.Unit)

		// Fetch transactions for specific account
		fmt.Printf("Press Enter to fetching transactions for account ID: %s...\n", accountId)
		fmt.Scanln()
		accountTransactions, err := comm.GetLastAccountTransactions(accountId)
		if err != nil {
			log.Println("Error fetching account transactions:", err)
			return
		}

		for _, transaction := range accountTransactions {
			fmt.Printf("Reference ID: %s, Amount: %s %s, Booking Date: %s, Infotext: %s\n", transaction.Reference, transaction.Amount.Value, transaction.Amount.Unit, transaction.BookingDate, transaction.RemittanceInfo)
		}
	}

	if depot {

		//Load depots
		fmt.Println("Press Enter to fetching depots...")
		fmt.Scanln()
		depots, err := comm.GetDepots()
		if err != nil {
			log.Println("Error fetching depots:", err)
			return
		}

		for _, depot := range depots {
			fmt.Printf("Depot ID: %s, Depot type: %s\n", depot.DepotId, depot.DepotType)
		}

		var depotId = depots[0].DepotId

		// Fetch depot details
		fmt.Printf("Press Enter to fetching details for depot ID: %s...\n", depotId)
		fmt.Scanln()
		depotDetailsResponse, err := comm.GetDepotPositions(depotId)
		if err != nil {
			log.Println("Error fetching depot details:", err)
			return
		}
		depotDetails := depotDetailsResponse.Values
		aggregated := depotDetailsResponse.Aggregated
		fmt.Printf("Current value: %s,  Profit/Loss: %s\n", aggregated.CurrentValue.Value, aggregated.ProfitLossPurchaseAbs)

		for _, position := range depotDetails {
			fmt.Printf("Position ID: %s, WKN: %s, Quantity: %s\n", position.PositionId, position.Wkn, position.Quantity.Value)
		}

		// Fetch the first depot position
		fmt.Printf("Press Enter to fetching details for position ID: %s...\n", depotDetails[0].PositionId)
		fmt.Scanln()
		positionDetails, err := comm.GetDepotPosition(depotId, depotDetails[0].PositionId)
		if err != nil {
			log.Println("Error fetching position details:", err)
			return
		}
		fmt.Printf("Position ID: %s, WKN: %s, Quantity: %s, Current Price: %s\n", positionDetails.PositionId, positionDetails.Wkn, positionDetails.Quantity.Value, positionDetails.CurrentPrice.Price.Value)

		// Fetch depot transactions
		fmt.Printf("Press Enter to fetching transactions for depot ID: %s...\n", depotId)
		fmt.Scanln()
		depotTransactions, err := comm.GetDepotTransactions(depotId)
		if err != nil {
			log.Println("Error fetching depot transactions:", err)
			return
		}

		for _, transaction := range depotTransactions {
			fmt.Printf("Transaction ID: %s, WKN: %s, Quantity: %s, Booking Date: %s\n", transaction.TransactionId, transaction.Instrument.Wkn, transaction.Quantity.Value, transaction.BookingDate)
		}
	}

	if instrument {

		// Fetch instrument details for TRAT0N
		fmt.Printf("Press Enter to fetching details for instrument with WKN: %s...\n", "TRAT0N")
		fmt.Scanln()
		instrumentDetails, err := comm.GetInstrumentDetails("TRAT0N")
		if err != nil {
			log.Println("Error fetching instrument details:", err)
			return
		}
		for _, instrument := range instrumentDetails {
			fmt.Printf("WKN: %s, ISIN: %s, Name: %s\n", instrument.Wkn, instrument.Isin, instrument.Name)
		}
	}

	if orders {

	}

	// Exit application
	fmt.Println("\nPress Enter to exit to  the app...")
	fmt.Scanln()
}
