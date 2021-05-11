/*
 SPDX-License-Identifier: Apache-2.0
*/

// ==== Invoke Logs ====
// peer chaincode invoke -C mychannel -n logs -c '{"Args":["init"]}'
// peer chaincode invoke -C mychannel -n logs -c '{"Args":["addLog","pk_DS","pk_DC","pk_DP","sk_data","status","operation"]}'

// ==== Query Logs ====
// peer chaincode query -C mychannel -n logs -c '{"Args":["readLog","pk_DS","pk_DC","pk_DP"]}'
// peer chaincode query -C mychannel -n logs -c '{"Args":["getHistoryForLog","pk_DS","pk_DC","pk_DP"]}'
// peer chaincode query -C mychannel -n logs -c '{"Args":["queryAll"]}'


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

type Log struct {
	ObjectType string	`json:"docType"` //docType is used to distinguish the various types of objects in state database
	SkData     string	`json:"SKData"`
	Timestamp  string	`json:"Timestamp"`
	Status     string	`json:"Status"`
	Operation  string	`json:"Operation"`
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
	fmt.Println("Log Init")
	return shim.Success(nil)
}

// ============================================================
// Invoke - Our entry point for Invocations
// ============================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "addLog" { // Add a log
		return t.addLog(stub, args)
	} else if function == "readLog" { // Read a log 
		return t.readLog(stub, args)
	} else if function == "getHistoryForLog" { // Get full history of a log
		return t.getHistoryForLog(stub, args)
	} else if function == "queryAll" { // Query all logs
		return t.queryAll(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// addLog - add a new Log, store into chaincode state
// ============================================================
func (t *SimpleChaincode) addLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6 (pk_DS, pk_DC, pk_DP, sk_data, status, operation)")
	}
	pk_DS		:= args[0]
	pk_DC		:= args[1]
	pk_DP		:= args[2]
	sk_data		:= args[3]
	status		:= args[4]
	operation	:= args[5]

	fmt.Println("- start adding Log")

	// ==== Check if Log already exists ====
	LogAsBytes, err := stub.GetState(pk_DS + "-" + pk_DC + "-" + pk_DP)
	if err != nil {
		return shim.Error("Failed to get Log: " + err.Error())
	} else if LogAsBytes != nil {
		fmt.Println("This Log already exists, then we update the Log: " + pk_DS + "-" + pk_DC + "-" + pk_DP)
		LogToUpdate := Log{}
		err = json.Unmarshal(LogAsBytes, &LogToUpdate) // Unmarshal it aka JSON.parse()
		if err != nil {
			return shim.Error(err.Error())
		}

		// Update timestamp and operation
		LogToUpdate.Timestamp = time.Now().UTC().String()
		LogToUpdate.Operation = operation
		LogToUpdateJSONasBytes, err := json.Marshal(LogToUpdate)
		if err != nil {
			return shim.Error(err.Error())
		}

		// === Update Log to state ===
		err = stub.PutState(pk_DS + "-" + pk_DC + "-" + pk_DP, LogToUpdateJSONasBytes)
		if err != nil {
			return shim.Error(err.Error())
		}

		fmt.Println("- end update Log")
		return shim.Success(nil)
	}

	// ==== Create Log object and marshal to JSON ====
	objectType := "Log"
	Log := &Log{objectType, sk_data, time.Now().UTC().String(), status, operation}
	LogJSONasBytes, err := json.Marshal(Log)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save Log to state ===
	err = stub.PutState(pk_DS + "-" + pk_DC + "-" + pk_DP, LogJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end upload Log")
	return shim.Success(nil)
}

// ============================================================
// readLog - read a Log from chaincode state
// ============================================================
func (t *SimpleChaincode) readLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var jsonResp string
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3 (pk_DS, pk_DC, pk_DP)")
	}

	pk_DS := args[0]
	pk_DC := args[1]
	pk_DP := args[2]

	valAsbytes, err := stub.GetState(pk_DS + "-" + pk_DC + "-" + pk_DP) // Get the Log from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + pk_DS + "-" + pk_DC + "-" + pk_DP + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Log does not exist: " + pk_DS + "-" + pk_DC + "-" + pk_DP + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// ============================================================
// getHistoryForLog - get the full history for a Log
// ============================================================
func (t *SimpleChaincode) getHistoryForLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3 (pk_DS, pk_DC, pk_DP)")
	}
	pk_DS := args[0]
	pk_DC := args[1]
	pk_DP := args[2]
	fmt.Printf("- start getHistoryForlog: %s\n", pk_DS + "-" + pk_DC + "-" + pk_DP)

	resultsIterator, err := stub.GetHistoryForKey(pk_DS + "-" + pk_DC + "-" + pk_DP)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the Consent
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
		//as-is (as the Value itself a JSON Consent)
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

	fmt.Printf("- getHistoryForLog returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// ============================================================
// queryAll - querayAll Logs
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
