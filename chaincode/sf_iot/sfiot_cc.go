/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"fmt"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Smartmeter struct {
	ObjectType string `json:"docType"`
	Id	string `json:"Id"`
        Timestamp string `json:"Timestamp"`
	Status  string `json:"Status"`
	Consumption string `json:"Consumption"`
        GatewayId string `json:GatewayId`
        GIdTimestamp string `json:"GIdTimestamp"`
}

type SmartContract struct {

}

var consumption = 0.0	// count for cosumption value
var abnormal = 0	//count for the number of abnormal data


func (t *SmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Smartmeter Init")
	return shim.Success(nil)
}

func (t *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Smartmeter Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "save_data" {
		return t.save_data(stub, args)
	} else if function == "delete_by_meterId" {
		return t.delete_by_meterId(stub, args)
	} else if function == "query_all" {
		return t.query_all(stub, args)
	} else if function == "query_by_meterId" {
                return t.query_by_meterId(stub, args)
	} else if function == "query_by_key" {
		return t.query_by_key(stub, args)
	} else if function == "query_by_month" {
                return t.query_by_month(stub, args, 0)
	} else if function == "compute_abnormal_by_month" {
		return t.compute_abnormal_by_month(stub, args)
	} else if function == "compute_consumption_by_month" {
                return t.compute_consumption_by_month(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"save_data\" \"delete_by_meterId\" \"query_all\" \"query_by_meterId\" \"query_by_key\" \"query_by_month\" \"compute_abnormal_by_month\" \"compute_consumption_by_month\"")
}

func (t *SmartContract) save_data(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5(Id, Timestamp, Status, Consumption, GatewayId)")
	}

	//DLA detection
	//if args[3] == "-999" {
	//	alarm()	
	//}

	//Create meter object and marshal to json (ex.)
	key := args[4] + "," + args[0] + "," + args[1]
	meter := &Smartmeter{"Smartmeter", args[0], args[1], args[2], args[3], args[4], key}


        // Check if meter already exists
        meterAsBytes, err := stub.GetState(key)

        if err != nil {
                return shim.Error("Failed to get meter: " + err.Error())
        } else if meterAsBytes != nil {

                // Retrieve the valAsbytes and check the input gatewayId
                Sep_status := Smartmeter{}
                err = json.Unmarshal(meterAsBytes, &Sep_status) //unmarshal it
                if err != nil {
                        return shim.Error(err.Error())
                }

		//Update the status and pusState
                Sep_status.Status = args[2]
		meterJsonasBytes, err := json.Marshal(Sep_status)
	        err = stub.PutState(key, meterJsonasBytes)
		if err != nil {
			return shim.Error(err.Error())
		}

		response := fmt.Sprintf("Update %s,%s status successful!", args[4], args[0])
		return  shim.Success([]byte(response))
	}


	meterasBytes, err := json.Marshal(meter)
        if err != nil {
                return shim.Error("json data error: " + err.Error())
        }

	err = stub.PutState(key, meterasBytes)
        if err != nil {
                return shim.Error("Save data error: " + err.Error())
        }

	//index the smartmeter to enable meter id-based query
        indexName := "Id~GIdTimestamp"
        idNameIndexkey, err := stub.CreateCompositeKey(indexName, []string{meter.Id, meter.GIdTimestamp})

        if err != nil {
                return shim.Error("Fail to create Composite key!")
}

	//Save index entry to state
        value := []byte{0x00}
        stub.PutState(idNameIndexkey, value)

	return shim.Success([]byte("Store data success!"))
}

// Deletes an entity from state
func (t *SmartContract) delete_by_meterId(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 2 {
                return shim.Error("Incorrect number of arguments. Expecting 2(Id, gatewayId)")
        }

        id := args[0]
        gatewayid := args[1]

	//Get State of CompositeKey "Id~GIdTimestamp"
	idResultsIterator, err := stub.GetStateByPartialCompositeKey("Id~GIdTimestamp", []string{id})

        if err != nil {
                return shim.Error("Fail to partial Composite key!")
        }
        defer idResultsIterator.Close()


        for idResultsIterator.HasNext() {
                queryResponse, err := idResultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }

                objectType, compositeKeys, err := stub.SplitCompositeKey(string(queryResponse.Key))
                fmt.Printf("%s", objectType)
                returnKey := compositeKeys[1]

                //fmt.Println(compositeKeys[0])	Id~GIdTimestampId
                //fmt.Println(compositeKeys[1])	GatewayId,Id,Timestamp

                // Get the attribute of return key
                valAsbytes, err := stub.GetState(returnKey)
                if err != nil {
                        return shim.Error("Failed to get smartmeter:" + err.Error())
                } else if valAsbytes == nil {
                        return shim.Error("Smartmeter does not exist")
                }

                // Retrieve the valAsbytes and check the input gatewayId
                Sep_gatewayId := Smartmeter{}
                err = json.Unmarshal(valAsbytes, &Sep_gatewayId) //unmarshal it
                if err != nil {
                        return shim.Error(err.Error())
                }

                if Sep_gatewayId.GatewayId == gatewayid {
			err := stub.DelState(returnKey)
		        if err != nil {
				return shim.Error("Failed to delete state")
			}
                }
	}

	return shim.Success([]byte("Delete data success!"))

}



func (s *SmartContract) query_by_key(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 3 {
                return shim.Error("Incorrect number of arguments. Expecting 3(Id, Timestamp, GatewaId)")
        }

	key := args[2] + "," + args[0] + "," + args[1]
        meterAsBytes, _ := stub.GetState(key)
        return shim.Success(meterAsBytes)
}


// query callback representing the query of a chaincode
func (t *SmartContract) query_all(stub shim.ChaincodeStubInterface, args []string) pb.Response {

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



func (t *SmartContract) query_by_meterId(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 2 {
                return shim.Error("Incorrect number of arguments. Expecting 2(Id, gatewayId)")
        }

        id := args[0]
        gatewayid := args[1]

	//Get State of CompositeKey "Id~GIdTimestamp"
	idResultsIterator, err := stub.GetStateByPartialCompositeKey("Id~GIdTimestamp", []string{id})

        if err != nil {
                return shim.Error("Fail to partial Composite key!")
        }
        defer idResultsIterator.Close()


        var buffer bytes.Buffer
        buffer.WriteString("\n[")

        bArrayMemberAlreadyWritten := false
        for idResultsIterator.HasNext() {
                queryResponse, err := idResultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }

                objectType, compositeKeys, err := stub.SplitCompositeKey(string(queryResponse.Key))
                fmt.Printf("%s", objectType)
                returnKey := compositeKeys[1]

                //fmt.Println(compositeKeys[0])	Id~GIdTimestampId
                //fmt.Println(compositeKeys[1])	GatewayId,Id,Timestamp

                // Get the attribute of return key
                valAsbytes, err := stub.GetState(returnKey)
                if err != nil {
                        return shim.Error("Failed to get smartmeter:" + err.Error())
                } else if valAsbytes == nil {
                        return shim.Error("Smartmeter does not exist")
                }

                // Retrieve the valAsbytes and check the input gatewayId
                Sep_gatewayId := Smartmeter{}
                err = json.Unmarshal(valAsbytes, &Sep_gatewayId) //unmarshal it
                if err != nil {
                        return shim.Error(err.Error())
                }

                if Sep_gatewayId.GatewayId == gatewayid {
                        if bArrayMemberAlreadyWritten == true {
                                buffer.WriteString(",")
                        }
                        buffer.WriteString("{\"Key\":")
                        buffer.WriteString("\"")
                        buffer.WriteString(returnKey)
                        buffer.WriteString("\"")
                        buffer.WriteString(", \"Record\":")
                        buffer.WriteString(string(valAsbytes))
                        buffer.WriteString("}")
                        bArrayMemberAlreadyWritten = true
                }
        }

        buffer.WriteString("]\n")
        return shim.Success(buffer.Bytes())

}


// flag used to distinquish query or compute function by month
func (t *SmartContract) query_by_month(stub shim.ChaincodeStubInterface, args []string, flag int) pb.Response {

        if len(args) != 4 {
                return shim.Error("Incorrect number of arguments. Expecting 4(Id, gatewayId, monthStart, monthEnd)")
        }

        id := args[0]
        gatewayid := args[1]
	monthStart := args[2]
	monthEnd := args[3]

	//Get State of CompositeKey "Id~GIdTimestamp"
	idResultsIterator, err := stub.GetStateByPartialCompositeKey("Id~GIdTimestamp", []string{id})

        if err != nil {
                return shim.Error("Fail to partial Composite key!")
        }
        defer idResultsIterator.Close()


        var buffer bytes.Buffer
        buffer.WriteString("\n[")

        bArrayMemberAlreadyWritten := false
        for idResultsIterator.HasNext() {
                queryResponse, err := idResultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }

                objectType, compositeKeys, err := stub.SplitCompositeKey(string(queryResponse.Key))
                fmt.Printf("%s", objectType)
                returnKey := compositeKeys[1]

                //fmt.Println(compositeKeys[0])	Id~GIdTimestampId
                //fmt.Println(compositeKeys[1])	GatewayId,Id,Timestamp

                // Get the attribute of return key
                valAsbytes, err := stub.GetState(returnKey)
                if err != nil {
                        return shim.Error("Failed to get smartmeter:" + err.Error())
                } else if valAsbytes == nil {
                        return shim.Error("Smartmeter does not exist")
                }

                // Retrieve the valAsbytes and check the input gatewayId
                Sep_gatewayId := Smartmeter{}
                err = json.Unmarshal(valAsbytes, &Sep_gatewayId) //unmarshal it
                if err != nil {
                        return shim.Error(err.Error())
                }

		// Compare GatewatId and Timestamp, >= monthStart && < (month of monthEnd)+1
                if Sep_gatewayId.GatewayId == gatewayid {
			monthOfEnd, err := strconv.Atoi(strings.SplitAfter(monthEnd, "-")[1])
			if err != nil {
				return shim.Error(err.Error())
			}

			if Sep_gatewayId.Timestamp >= monthStart && Sep_gatewayId.Timestamp < (strings.SplitAfter(monthEnd, "-")[0] + strconv.Itoa(monthOfEnd + 1)) {
				if flag == 0 {

					if bArrayMemberAlreadyWritten == true {
						buffer.WriteString(",")
					}
					buffer.WriteString("{\"Key\":")
					buffer.WriteString("\"")
					buffer.WriteString(returnKey)
					buffer.WriteString("\"")
					buffer.WriteString(", \"Record\":")
					buffer.WriteString(string(valAsbytes))
					buffer.WriteString("}")
					bArrayMemberAlreadyWritten = true

				}else if Sep_gatewayId.Status == "A"{
					abnormal += 1
				}

				// Change string to float and compute the consumption
				parse, err := strconv.ParseFloat(Sep_gatewayId.Consumption, 8)
				if err != nil {
		                        return shim.Error("Failed to compute the consumption:")
				}
				consumption += parse
			}
		}
        }

        buffer.WriteString("]\n")
        return shim.Success(buffer.Bytes())
}



func (t *SmartContract) compute_abnormal_by_month(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 4 {
                return shim.Error("Incorrect number of arguments. Expecting 4(Id, gatewayId, monthStart, monthEnd)")
        }

        abnormal, consumption = 0, 0.0
	t.query_by_month(stub, args, 1)
        response := fmt.Sprintf("The number of abnormal datas is %d !", abnormal)

        return  shim.Success([]byte(response))
}



func (t *SmartContract) compute_consumption_by_month(stub shim.ChaincodeStubInterface, args []string) pb.Response {

        if len(args) != 4 {
                return shim.Error("Incorrect number of arguments. Expecting 4(Id, gatewayId, monthStart, monthEnd)")
        }

        abnormal, consumption = 0, 0.0
        t.query_by_month(stub, args, 1)
        response := fmt.Sprintf("The total consumption is %f !", consumption)

        return  shim.Success([]byte(response))
}



func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error starting Smart Contract: %s", err)
	}
}

