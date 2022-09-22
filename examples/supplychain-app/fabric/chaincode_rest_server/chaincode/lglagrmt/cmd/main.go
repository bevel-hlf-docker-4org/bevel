package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	lglagrmt "github.com/chaincode/lglagrmt"
)

// main function starts up the chaincode in the container during instantiate
func main() {
	err := shim.Start(new(lglagrmt.SmartContract))

	if err != nil {
		fmt.Printf("Error starting lglagrmt chaincode: %s", err)
	}
}
