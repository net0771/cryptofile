/*
 * Copyright IBM Corp All Rights Reserved
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"bytes"
	"time"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// SimpleAsset implements a simple chaincode to manage an asset
type SimpleAsset struct {
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the args from the transaction proposal
	args := stub.GetStringArgs()
	if len(args) != 2 {
		return shim.Error("Incorrect arguments. Expecting a key and a value")
	}

	// Set up any variables or assets here by calling stub.PutState()

	// We store the key and the value on the ledger
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	}
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error
	if fn == "set" {
		result, err = set(stub, args)
	} else if fn == "transfer" {
		result, err = transfer(stub, args)
	} else if fn == "delete" {
		result, err = delete(stub, args)
	} else { // assume 'get' even if fn is nil
		result, err = get(stub, args)
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}

// Set stores the asset (both key and value) on the ledger. If the key exists,
// it will override the value with the new one
func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return args[1], nil
}

// Get returns the value of the specified asset key
func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	value, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
		return "", fmt.Errorf("Asset not found: %s", args[0])
	}
	
	var buffer bytes.Buffer
	buffer.WriteString("[")
	resultsIterator, err := stub.GetHistoryForKey(args[0])
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	defer resultsIterator.Close()
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return "", fmt.Errorf(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForMarble returning:\n%s\n", buffer.String())

	return string(value), nil
}

// Transfer A to B
func transfer(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 3 { // ["user1", "user2", "100"]
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	var To, From string
	var ToValue, FromValue int
	var TransferValue int
	var err error

	From = args[0]
	To = args[1]
	FromBytesVal, err := stub.GetState(From)
	if err != nil {
		return "", fmt.Errorf("Failed to get state")
	}
	if FromBytesVal == nil {
		return "", fmt.Errorf("The key not found")
	}

	ToBytesVal, _ := stub.GetState(To)
	if err != nil {
		return "", fmt.Errorf("Failed to get state")
	}
	if ToBytesVal == nil {
		return "", fmt.Errorf("The key not found")
	}
	
	FromValue, _ = strconv.Atoi(string(FromBytesVal)) // []bytes -> string  -> int
	ToValue, _ = strconv.Atoi(string(ToBytesVal)) // []bytes -> string  -> int
	TransferValue, err = strconv.Atoi(args[2])
	if err != nil {
		return "", fmt.Errorf("Invalid transaction amount, please use a integer value")
	}

	if FromValue < TransferValue {
		return "", fmt.Errorf("Insufficient amount")
	}

	// Transfer 
	FromValue = FromValue - TransferValue	
	ToValue = ToValue + TransferValue
	fmt.Printf("From Value = %d, To Value = %d", FromValue, ToValue)

	// PutState
	err = stub.PutState(From, []byte(strconv.Itoa(FromValue))) // int -> string -> []byte
	if err != nil {
		return "", fmt.Errorf("Failed to transfer")
	}
	err = stub.PutState(To, []byte(strconv.Itoa(ToValue))) // int -> string -> []byte
	if err != nil {
		return "", fmt.Errorf("Failed to transfer")
	}

	return args[2], nil
}

// Delete key from the World State(State DB)
func delete(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	err := stub.DelState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to delete state")
	} else {
		value, _ := stub.GetState(args[0])
		fmt.Printf("Deleted:" + string(value))
	}
	return args[0] + " was deleted.", nil
}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}
