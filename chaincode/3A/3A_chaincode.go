/*
 SPDX-License-Identifier: Apache-2.0
*/

// ==== Invoke Consents ====
// peer chaincode invoke -C mychannel -n consents -c '{"Args":["init"]}'
// peer chaincode invoke -C mychannel -n consents -c '{"Args":["grantConsent","pk_DS","pk_DC","policy","enhash","pk_enc"]}'

// ==== Query Consents ====
// peer chaincode query -C mychannel -n consents -c '{"Args":["getHistoryForConsent","pk_DS","pk_DC"]}'
// peer chaincode query -C mychannel -n consents -c '{"Args":["readConsent","pk_DS","pk_DC"]}'
// peer chaincode query -C mychannel -n consents -c '{"Args":["queryAll"]}'


package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"io/ioutil"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Consent struct {
	ObjectType string				`json:"docType"` //docType is used to distinguish the various types of objects in state database
	Policy	   map[string][]string	`json:"Policy"`
	Timestamp  string				`json:"Timestamp"`
	Enhash	   string 				`json:"Enhash"`
	pk_enc	   string 				`json:"pk_enc"`
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
	fmt.Println("Consent Init")
	return shim.Success(nil)
}

// ============================================================
// Invoke - Our entry point for Invocations
// ============================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "grantConsent" { // Grant consent 
		return t.grantConsent(stub, args)
	} else if function == "readConsent" { // Read a consent 
		return t.readConsent(stub, args)
	} else if function == "getHistoryForConsent" { // Get full history of conesent record
		return t.getHistoryForConsent(stub, args)
	} else if function == "queryAll" { // Query all consents
		return t.queryAll(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) // error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// grantConsent - Grant Consent to DP (or AG)
// ============================================================
func (t *SimpleChaincode) grantConsent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5 (pk_DS, pk_DC, policy, enhash, pk_enc)")
	}
	pk_DS	:= args[0]
	pk_DC	:= args[1]
	policy	:= args[2]	// e.g., "C": "+pk_DP", "R": "-pk_DP", ...
	enhash	:= args[3]
	pk_enc	:= args[4]

	fmt.Println("- start adding Consent")

	// ==== Check if Consent already exists ====
	ConsentAsBytes, err := stub.GetState(pk_DS + "-" + pk_DC)
	if err != nil {
		return shim.Error("Failed to get Consent: " + err.Error())
	} else if ConsentAsBytes != nil {
		// ==== Update consent ====
		fmt.Println("This Consent already exists: " + pk_DS + "-" + pk_DC)

		// Split input policy (e.g., "C":"+pk_DP","R":"-pk_DP", ...)
		ConsentToUpdate := Consent{}
		err = json.Unmarshal(ConsentAsBytes, &ConsentToUpdate) // Unmarshal it aka JSON.parse()
		if err != nil {
			return shim.Error(err.Error())
		}
		
		singlePolicy := strings.Split(policy, ",")
		var index int

		for _, v := range singlePolicy{
			// Retrieve key and value of input policy
			v = strings.Trim(v, "{}")
			key := string(strings.Split(v, ":")[0][1])
			value := strings.Split(v, ":")[1]
			value = strings.Trim(value , "\"")
			policies := ConsentToUpdate.Policy[key]

			// Add(+) or Remove(-) consent
			if value[0] == '+' {
				policies = append(policies, value[1:])
			} else if value[0] == '-' {
				for i, v := range policies {
					if v == value[1:] {
						index = i
						break
					}
				}

				// Remove the consent whose address is var "index"
				policies[len(policies)-1], policies[index] = policies[index], policies[len(policies)-1]
				policies = policies[:len(policies)-1]	
				
			} else {
				return shim.Error("Invalid input consent list!")
			}
			
			// Update consent to update which is going to be written to the state
			ConsentToUpdate.Policy[key] = policies
			fmt.Println(ConsentToUpdate.Policy[key])
		}

		ConsentJSONasBytes, _ := json.Marshal(ConsentToUpdate)
		err = stub.PutState(pk_DS + "-" + pk_DC, ConsentJSONasBytes) //Update the pPolicy
		if err != nil {
			return shim.Error(err.Error())
		}

		fmt.Println("- end updating Consent (success)")
		return shim.Success(nil)
	} else {
		// ==== Add new consent ====
		// Read encrpyted hash first
		enhash, err := ioutil.ReadFile("./" + enhash)
		if err != nil {
			return shim.Error("Failed to get encrpyted hash: " + err.Error())
		}

		p := make(map[string][]string)
		p["C"] = append(p["C"], pk_DC)
		p["R"] = append(p["R"], pk_DC)
		p["U"] = append(p["U"], pk_DC)
		p["D"] = append(p["D"], pk_DC)

		objectType := "Consent"
		Consent := &Consent{objectType, p, time.Now().UTC().String(), string(enhash), pk_enc}
		ConsentJSONasBytes, err := json.Marshal(Consent)
		if err != nil {
			return shim.Error(err.Error())
		}
	
		// === Save Consent to state ===
		err = stub.PutState(pk_DS + "-" + pk_DC, ConsentJSONasBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
	
		fmt.Println("- end adding Consent (success)")
		return shim.Success(nil)
	}
}

// ============================================================
// getHistoryForConsent - get the full history for a Consent
// ============================================================
func (t *SimpleChaincode) getHistoryForConsent(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2 (pk_DS, pk_DC)")
	}
	pk_DS := args[0]
	pk_DC := args[1]
	fmt.Printf("- start getHistoryForConsent: %s\n", pk_DS + "-" + pk_DC)

	resultsIterator, err := stub.GetHistoryForKey(pk_DS + "-" + pk_DC)
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

	fmt.Printf("- getHistoryForConsent returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// ============================================================
// readConsent - read a Consent from chaincode state
// ============================================================
func (t *SimpleChaincode) readConsent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var jsonResp string
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2 (pk_DS, pk_DC)")
	}

	pk_DS := args[0]
	pk_DC := args[1]

	valAsbytes, err := stub.GetState(pk_DS + "-" + pk_DC) // Get the Consent from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + pk_DS + "-" + pk_DC + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Log does not exist: " + pk_DS + "-" + pk_DC + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// ============================================================
// queryAll - querayAll Consents
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
