package main

import (
	"fmt"

	"github.com/fritzrepo/stockportfolio/depot"
)

func main() {
	depot.ComputeTransactions()
	fmt.Println("End")
}
