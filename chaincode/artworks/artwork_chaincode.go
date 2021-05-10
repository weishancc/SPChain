/*
 SPDX-License-Identifier: Apache-2.0
*/

// ==== Invoke Artworks ====
// peer chaincode invoke -C mychannel -n artworks -c '{"Args":["init"]}'
// peer chaincode invoke -C mychannel -n artworks -c '{"Args":["uploadArtwork","tokenID","multiHash","owner","creator"]}'
// peer chaincode invoke -C mychannel -n artworks -c '{"Args":["transferArtwork","tokenID","newOwner","multiHash"]}'
// peer chaincode invoke -C mychannel -n artworks -c '{"Args":["deleteArtwork","tokenID"]}'

// ==== Query Artworks ====
// peer chaincode query -C mychannel -n artworks -c '{"Args":["readArtwork","tokenID"]}'
// peer chaincode query -C mychannel -n artworks -c '{"Args":["getHistoryForArtwork","tokenID"]}'
// peer chaincode query -C mychannel -n artworks -c '{"Args":["queryAll"]}'


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

type Artwork struct {
	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Multihash  string `json:"Multihash"`
	Owner      string `json:"Owner"`
	Creator    string `json:"Creator"`
	Timestamp  string `json:"Timestamp"`
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
	fmt.Println("Artwork Init")
	return shim.Success(nil)
}

// ============================================================
// Invoke - Our entry point for Invocations
// ============================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "uploadArtwork" { //create a new Artwork
		return t.uploadArtwork(stub, args)
	} else if function == "transferArtwork" { //change owner of a specific Artwork
		return t.transferArtwork(stub, args)
	} else if function == "deleteArtwork" { //delete a Artwork
		return t.deleteArtwork(stub, args)
	} else if function == "readArtwork" { //read a Artwork
		return t.readArtwork(stub, args)
	} else if function == "getHistoryForArtwork" { //get history of values for a Artwork
		return t.getHistoryForArtwork(stub, args)
	} else if function == "queryAll" { //get history of values for a Artwork
		return t.queryAll(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// initArtwork - create a new Artwork, store into chaincode state
// ============================================================
func (t *SimpleChaincode) uploadArtwork(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4 (TokenID, Multi-Hash, Owner, Creator)")
	}

	// ==== Input sanitation ====
	fmt.Println("- start uploading Artwork")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}

	tokenID     := args[0]
	multiHash	:= args[1]
	owner		:= args[2]
	creator		:= args[3]

	// ==== Check if Artwork already exists ====
	ArtworkAsBytes, err := stub.GetState(tokenID)
	if err != nil {
		return shim.Error("Failed to get Artwork: " + err.Error())
	} else if ArtworkAsBytes != nil {
		fmt.Println("This Artwork already exists: " + tokenID)
		return shim.Error("This Artwork already exists: " + tokenID)
	}

	// ==== Create Artwork object and marshal to JSON ====
	//Read encrpyted hash first
	//enhash, err := ioutil.ReadFile("./" + en_pointer)
	//if err != nil {
	//	return shim.Error("Failed to get encrpyted hash: " + err.Error())
	//}

	objectType := "Artwork"
	artwork := &Artwork{objectType, multiHash, owner, creator, time.Now().UTC().String()}
	ArtworkJSONasBytes, err := json.Marshal(artwork)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save Artwork to state ===
	err = stub.PutState(tokenID, ArtworkJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== Artwork uploaded a. Return success ====
	fmt.Println("- end upload Artwork")
	return shim.Success(nil)
}

// ============================================================
// readArtwork - read a Artwork from chaincode state
// ============================================================
func (t *SimpleChaincode) readArtwork(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var tokenID, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting TokenID of the Artwork to query")
	}

	tokenID = args[0]
	valAsbytes, err := stub.GetState(tokenID) //get the Artwork from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + tokenID + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Artwork does not exist: " + tokenID + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// ============================================================
// deleteArtwork - remove a Artwork key/value pair from state
// ============================================================
func (t *SimpleChaincode) deleteArtwork(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var jsonResp string
	var ArtworkJSON Artwork
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting TokenID.")
	}
	tokenID := args[0]

	valAsbytes, err := stub.GetState(tokenID) //get the Artwork from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + tokenID + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Artwork does not exist: " + tokenID + "\"}"
		return shim.Error(jsonResp)
	}

	err = json.Unmarshal([]byte(valAsbytes), &ArtworkJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + tokenID + "\"}"
		return shim.Error(jsonResp)
	}

	err = stub.DelState(tokenID) //remove the Artwork from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	return shim.Success(nil)
}

// ===============================================================================
// transferArtwork - transfer a Artwork by setting a new owner name on the Artwork
// ===============================================================================
func (t *SimpleChaincode) transferArtwork(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3 (TokenID, NewOwner, MultiHash)")
	}

	tokenID := args[0]
	newOwner := args[1]
	multiHash := args[2]
	fmt.Println("- start transferArtwork ", tokenID, newOwner)

	ArtworkAsBytes, err := stub.GetState(tokenID)
	if err != nil {
		return shim.Error("Failed to get Artwork:" + err.Error())
	} else if ArtworkAsBytes == nil {
		return shim.Error("Artwork does not exist")
	}

	ArtworkToTransfer := Artwork{}
	err = json.Unmarshal(ArtworkAsBytes, &ArtworkToTransfer) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	ArtworkToTransfer.Owner = newOwner //change the owner
	ArtworkToTransfer.Multihash =  multiHash //update multihash

	ArtworkJSONasBytes, _ := json.Marshal(ArtworkToTransfer)
	err = stub.PutState(tokenID, ArtworkJSONasBytes) //rewrite the Artwork
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end transferArtwork (success)")
	return shim.Success(nil)
}

// ============================================================
// getHistoryForArtwork - get the full history for a Artwork
// ============================================================
func (t *SimpleChaincode) getHistoryForArtwork(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1 (TokenID)")
	}

	tokenID := args[0]

	fmt.Printf("- start getHistoryForArtwork: %s\n", tokenID)

	resultsIterator, err := stub.GetHistoryForKey(tokenID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the Artwork
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
		//as-is (as the Value itself a JSON Artwork)
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

	fmt.Printf("- getHistoryForArtwork returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// ============================================================
// queryAll - querayAll Artworks
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
