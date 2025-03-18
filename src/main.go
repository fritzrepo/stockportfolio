package main

import (
	"fmt"

	"github.com/fritzrepo/stockportfolio/depot"
)

// contains prüft, ob ein Wert im Slice vorhanden ist
// func contains(slice []depotEntry, value string) bool {
// 	for _, v := range slice {
// 		if v.asset == value {
// 			return true
// 		}
// 	}
// 	return false
// }

func main() {
	depot.ComputeTransactions()
	fmt.Println("End")
}
