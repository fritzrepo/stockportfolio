package main

import (
	"fmt"
	"log"

	"github.com/fritzrepo/stockportfolio/internal/banking/comdirect"
)

func main() {

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

	fmt.Printf("Press Enter to fetching balance for account ID: %s...\n", accountId)
	fmt.Scanln()
	accountBalance, err := comm.GetAccoutBalance(accountId)
	if err != nil {
		log.Println("Error fetching account balance:", err)
		return
	}
	fmt.Printf("Account ID: %s, Balance: %s %s\n", accountBalance.Account.AccountType.Text, accountBalance.Balance.Value, accountBalance.Balance.Unit)

	fmt.Println("\nPress Enter to exit to  the app...")
	fmt.Scanln()
}
