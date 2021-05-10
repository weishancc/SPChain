/*
 SPDX-License-Identifier: Apache-2.0
*/

// ==== Invoke Wallets ====
// peer chaincode invoke -C mychannel -n wallets -c '{"Args":["init"]}'
// peer chaincode invoke -C mychannel -n wallets -c '{"Args":["addWallet","name","price"]}'
// peer chaincode invoke -C mychannel -n wallets -c '{"Args":["transferBalance","payer","payee","creator","price","r_y"]}'
// peer chaincode invoke -C mychannel -n wallets -c '{"Args":["deleteWallet","name"]}'

// ==== Query Wallets ====
// peer chaincode query -C mychannel -n wallets -c '{"Args":["readBalance","name"]}'
// peer chaincode query -C mychannel -n wallets -c '{"Args":["getHistoryForWallet","name"]}'
// peer chaincode query -C mychannel -n wallets -c '{"Args":["queryAll"]}'


package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Wallet struct {
	ObjectType string	`json:"docType"` //docType is used to distinguish the various types of objects in state database
	Balance    float64	`json:"Balance"`
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ============================================================
// Init initializes chaincode
// ============================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Wallet Init")
	return shim.Success(nil)
}

// ============================================================
// Invoke - Our entry point for Invocations
// ============================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "addWallet" { //Add a Wallet
		return t.addWallet(stub, args)
	} else if function == "readBalance" { //read a Wallet balance
		return t.readBalance(stub, args)
	} else if function == "transferBalance" { //transfer balance
		return t.transferBalance(stub, args)
	} else if function == "deleteWallet" { //delete a Wallet
		return t.deleteWallet(stub, args)
	} else if function == "getHistoryForWallet" { //get full history of a Wallet (transfer record)
		return t.getHistoryForWallet(stub, args)
	} else if function == "queryAll" { //query all Wallets
		return t.queryAll(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// addWallet - add a new Wallet, store into chaincode state
// ============================================================
func (t *SimpleChaincode) addWallet(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2 (Name, Balance)")
	}
	name := args[0]
	balance := args[1]
	
	fmt.Println("- start adding Wallet")

	// ==== Check if Wallet already exists ====
	WalletAsBytes, err := stub.GetState(name)
	if err != nil {
		return shim.Error("Failed to get Wallet: " + err.Error())
	} else if WalletAsBytes != nil {
		fmt.Println("This Wallet already exists: " + name)
		return shim.Error("This Wallet already exists: " + name)
	}

	// ==== Create Wallet object and marshal to JSON ====
	objectType := "Wallet"
	fbalance, _ := strconv.ParseFloat(balance, 64)	// Convert balance string to float64 type
	wallet := &Wallet{objectType, fbalance}
	WalletJSONasBytes, err := json.Marshal(wallet)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save Wallet to state ===
	err = stub.PutState(name, WalletJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end upload Wallet")
	return shim.Success(nil)
}

// ============================================================
// readBalance - read a Wallet balance from chaincode state
// ============================================================
func (t *SimpleChaincode) readBalance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the Wallet to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name) //get the Wallet from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Wallet does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// ===============================================================================
// transferBalance - transfer a Wallet by setting a new owner name on the Wallet
// ===============================================================================
func (t *SimpleChaincode) transferBalance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5 (Payer, Payee, Creator, Price, R_y)")
	}

	payer 	:= args[0]
	payee	:= args[1]
	creator := args[2]
	price	:= args[3]
	r_y		:= args[4]
	fmt.Println("- start transferBalance ", payer, payee, price)

	// Payer (Collector)
	PayerWalletAsBytes, err := stub.GetState(payer)
	if err != nil {
		return shim.Error("Failed to get Wallet:" + err.Error())
	} else if PayerWalletAsBytes == nil {
		return shim.Error("Wallet does not exist")
	}

	PayerWalletToTransfer := Wallet{}
	err = json.Unmarshal(PayerWalletAsBytes, &PayerWalletToTransfer) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}

	fr_y, _ := strconv.ParseFloat(r_y, 64)
	fbalance, _ := strconv.ParseFloat(price, 64)
	PayerWalletToTransfer.Balance = PayerWalletToTransfer.Balance - fbalance //subtract the price

	PayerWalletJSONasBytes, _ := json.Marshal(PayerWalletToTransfer)
	err = stub.PutState(payer, PayerWalletJSONasBytes) //rewrite the Wallet
	if err != nil {
		return shim.Error(err.Error())
	}

	// Payee (Owner)
	PayeeWalletAsBytes, err := stub.GetState(payee)
	if err != nil {
		return shim.Error("Failed to get Wallet:" + err.Error())
	} else if PayeeWalletAsBytes == nil {
		return shim.Error("Wallet does not exist")
	}

	PayeeWalletToTransfer := Wallet{}
	err = json.Unmarshal(PayeeWalletAsBytes, &PayeeWalletToTransfer) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}

	// Payee (Owner) will not get royalties if payee euqals to creator (i.e., firsthand transaction)
	if payee == creator {
		PayeeWalletToTransfer.Balance = PayeeWalletToTransfer.Balance + fbalance //add the price
		PayeeWalletJSONasBytes, _ := json.Marshal(PayeeWalletToTransfer)
		err = stub.PutState(payee, PayeeWalletJSONasBytes) //rewrite the Wallet
		if err != nil {
			return shim.Error(err.Error())
		}
	} else {
		// Payee (Owner) get [price*(1-r_y)]
		PayeeWalletToTransfer.Balance = PayeeWalletToTransfer.Balance + fbalance * (1.0 - fr_y) //add the price
		PayeeWalletJSONasBytes, _ := json.Marshal(PayeeWalletToTransfer)
		err = stub.PutState(payee, PayeeWalletJSONasBytes) //rewrite the Wallet
		if err != nil {
			return shim.Error(err.Error())
		}

		// Orginal creator get royalties [price*r_y]
		CreatorWalletAsBytes, err := stub.GetState(creator)
		if err != nil {
			return shim.Error("Failed to get Wallet:" + err.Error())
		} else if CreatorWalletAsBytes == nil {
			return shim.Error("Wallet does not exist")
		}

		CreatorWalletToTransfer := Wallet{}
		err = json.Unmarshal(CreatorWalletAsBytes, &CreatorWalletToTransfer) //unmarshal it aka JSON.parse()
		if err != nil {
			return shim.Error(err.Error())
		}

		CreatorWalletToTransfer.Balance = CreatorWalletToTransfer.Balance + fbalance * fr_y
		CreatorWalletJSONasBytes, _ := json.Marshal(CreatorWalletToTransfer)
		err = stub.PutState(creator, CreatorWalletJSONasBytes) //rewrite the Wallet
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	fmt.Println("- end transferBalance (success)")
	return shim.Success(nil)
}

// ============================================================
// deleteWallet - remove a Wallet key/value pair from state
// ============================================================
func (t *SimpleChaincode) deleteWallet(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var jsonResp string
	var WalletJSON Wallet
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name.")
	}
	name := args[0]

	valAsbytes, err := stub.GetState(name) //get the Wallet from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Wallet does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	err = json.Unmarshal([]byte(valAsbytes), &WalletJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	err = stub.DelState(name) //remove the Wallet from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	return shim.Success(nil)
}

// ============================================================
// getHistoryForWallet - get the full history for a Wallet
// ============================================================
func (t *SimpleChaincode) getHistoryForWallet(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting Name.")
	}

	name := args[0]

	fmt.Printf("- start getHistoryForWallet: %s\n", name)

	resultsIterator, err := stub.GetHistoryForKey(name)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the Wallet
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
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
		//as-is (as the Value itself a JSON Wallet)
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

	fmt.Printf("- getHistoryForWallet returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// ============================================================
// queryAll - querayAll Wallets
// ============================================================
func (t *SimpleChaincode) queryAll(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Expecting 0")
	}

        resultsIterator, err := stub.GetStateByRange("","")
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()

        // buffer is a JSON array containing QueryResults
        var buffer bytes.Buffer
        buffer.WriteString("\n[")

	bArrayMemberAlreadyWritten := false
        for resultsIterator.HasNext() {
                queryResponse, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                // Add a comma before array members, suppress it for the first array member
                if bArrayMemberAlreadyWritten == true {
                        buffer.WriteString(",")
                }
                buffer.WriteString("{\"Key\":")
                buffer.WriteString("\"")
                buffer.WriteString(queryResponse.Key)
                buffer.WriteString("\"")

                buffer.WriteString(", \"Record\":")
                // Record is a JSON object, so we write as-is
                buffer.WriteString(string(queryResponse.Value))
                buffer.WriteString("}")
                bArrayMemberAlreadyWritten = true
        }
        buffer.WriteString("]\n")
        return shim.Success(buffer.Bytes())
}
