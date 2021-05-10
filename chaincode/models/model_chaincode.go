/*
 SPDX-License-Identifier: Apache-2.0
*/

// ==== Invoke Model ====
// peer chaincode invoke -C mychannel -n models -c '{"Args":["init"]}'
// peer chaincode invoke -C mychannel -n models -c '{"Args":["addModel","name","address","creator","desc"]}'
// peer chaincode invoke -C mychannel -n models -c '{"Args":["invokeModel","api","offset"]}'
// peer chaincode invoke -C mychannel -n models -c '{"Args":["deleteModel","name"]}'

// ==== Query Model ====
// peer chaincode query -C mychannel -n models -c '{"Args":["readModel","name"]}'
// peer chaincode query -C mychannel -n models -c '{"Args":["queryAll"]}'


package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
	"io/ioutil"
	"net/http"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Model struct {
	ObjectType string	`json:"docType"` //docType is used to distinguish the various types of objects in state database
	Address    string	`json:"Address"`
	Creator    string	`json:"Creator"`
	Desc       string	`json:"Desc"`
	Timestamp  string 	`json:"Timestamp"`
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
	fmt.Println("Model Init")
	return shim.Success(nil)
}

// ============================================================
// Invoke - Our entry point for Invocations
// ============================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "addModel" { //Add a Model
		return t.addModel(stub, args)
	} else if function == "readModel" {	//Read Model information
		return t.readModel(stub, args)
	} else if function == "deleteModel" { //Delete a Model
		return t.deleteModel(stub, args)
	} else if function == "invokeModel" { //Invoke Model API
		return t.invokeModel(stub, args)
	} else if function == "queryAll" { //Query all Model
		return t.queryAll(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// addModel - add a new Model, store into chaincode state
// ============================================================
func (t *SimpleChaincode) addModel(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4 (Name, Hash, Creator, Desc)")
	}
	name	:= args[0]
	address	:= args[1]
	creator	:= args[2]
	desc	:= args[3]

	fmt.Println("- start adding Model")

	// ==== Check if Model already exists ====
	ModelAsBytes, err := stub.GetState(name)
	if err != nil {
		return shim.Error("Failed to get Model: " + err.Error())
	} else if ModelAsBytes != nil {
		fmt.Println("This Model already exists: " + name)
		return shim.Error("This Model already exists: " + name)
	}

	// ==== Create Model object and marshal to JSON ====
	objectType := "Model"
	Model := &Model{objectType, address, creator, desc, time.Now().UTC().String()}
	ModelJSONasBytes, err := json.Marshal(Model)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save Model to state ===
	err = stub.PutState(name, ModelJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end upload Model")
	return shim.Success(nil)
}

// ============================================================
// readModel - read a Model from chaincode state
// ============================================================
func (t *SimpleChaincode) readModel(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the Model to query")
	}

	name := args[0]
	valAsbytes, err := stub.GetState(name) //get the Model from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Model does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// ============================================================
// deleteModel - remove a Model key/value pair from state
// ============================================================
func (t *SimpleChaincode) deleteModel(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var jsonResp string
	var ModelJSON Model

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name.")
	}
	name := args[0]

	valAsbytes, err := stub.GetState(name) //get the Model from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Model does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	err = json.Unmarshal([]byte(valAsbytes), &ModelJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	err = stub.DelState(name) //remove the Model from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	return shim.Success(nil)
}

// ============================================================
// invokeModel - invoke Model API to inference
// ============================================================
func (t *SimpleChaincode) invokeModel(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2 (api, offset)")
	}

	api := args[0]
	offset := args[1]

	m := make(map[string]string)
	m["image"] = offset
	jsonStr, _ := json.Marshal(m)

	// Post service
	resp, err := http.Post(api, "application/json", bytes.NewReader(jsonStr))
	if err != nil {
		return shim.Error(err.Error())
	}

	// Get the response
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	return shim.Success(nil)
}

// ============================================================
// queryAll - querayAll Model
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

        // buffer is a JSON array containing Query Results
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
