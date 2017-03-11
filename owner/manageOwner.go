/*/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at
  http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
"errors"
"fmt"
"strconv"
"encoding/json"

"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ManageOwner example simple Chaincode implementation
type ManageOwner struct {
}

var MerchantIndexStr = "_Merchantindex"				//name for the key/value that will store a list of all known Merchants

type Merchant struct{							// Attributes of a Merchant
	MerchantID string `json:"merchantId"`					
	MerchantName string `json:"merchantName"`
	MerchantLogo string `json:"merchantLogo"`
	MerchantIndustry string `json:"merchantIndustry"`					
	PointsPerDollarSpent string `json:"pointsPerDollarSpent"`
	MerchantCurrency string `json:"merchantCurrency"`
	ExchangeRate string `json:"exchangeRate"`
	MerchantCU_date string `json:"merchantCU_date"`
}

// ============================================================================================================================
// Main - start the chaincode for Merchant management
// ============================================================================================================================
func main() {			
	err := shim.Start(new(ManageOwner))
	if err != nil {
		fmt.Printf("Error starting Merchant management chaincode: %s", err)
	}
}
// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *ManageOwner) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var msg string
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting ' ' as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// Initialize the chaincode
	msg = args[0]
	// Write the state to the ledger
	err = stub.PutState("abc", []byte(msg))				//making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(MerchantIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	tosend := "{ \"message\" : \"ManageOwner chaincode is deployed successfully.\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 
	return nil, nil
}
// ============================================================================================================================
// Run - Our entry point for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func (t *ManageOwner) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}
// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *ManageOwner) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "createMerchant" {											//create a new Merchant
		return t.createMerchant(stub, args)
	}else if function == "deleteMerchant" {									// delete a Merchant
		return t.deleteMerchant(stub, args)
	}else if function == "updateMerchant" {									//update a Merchant
		return t.updateMerchant(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)
	errMsg := "{ \"message\" : \"Received unknown function invocation\", \"code\" : \"503\"}"
	err := stub.SetEvent("errEvent", []byte(errMsg))
	if err != nil {
		return nil, err
	} 
	return nil, nil			//error
}
// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *ManageOwner) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "getMerchantByID" {													//Read a Merchant by merchantId
		return t.getMerchantByID(stub, args)
	} else if function == "getAllMerchants" {													//Read all Merchants
		return t.getAllMerchants(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error
	errMsg := "{ \"message\" : \"Received unknown function query\", \"code\" : \"503\"}"
	err := stub.SetEvent("errEvent", []byte(errMsg))
	if err != nil {
		return nil, err
	} 
	return nil, nil
}
// ============================================================================================================================
// getMerchantByID - get Merchant details for a specific ID from chaincode state
// ============================================================================================================================
func (t *ManageOwner) getMerchantByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var merchantId string
	var err error
	fmt.Println("start getMerchantByID")
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 'merchantId' as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set merchantId
	merchantId = args[0]
	valAsbytes, err := stub.GetState(merchantId)									//get the merchantId from chaincode state
	if err != nil {
		errMsg := "{ \"message\" : \""+ merchantId + " not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Println("end getMerchantByID")
	return valAsbytes, nil													//send it onward
}
// ============================================================================================================================
//  getAllMerchants- get details of all Merchants from chaincode state
// ============================================================================================================================
func (t *ManageOwner) getAllMerchants(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, errResp string
	var merchantIndex []string
	var err error
	fmt.Println("start getAllMerchants")
		
	merchantAsBytes, err := stub.GetState(MerchantIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Merchant index")
	}
	json.Unmarshal(merchantAsBytes, &merchantIndex)								//un stringify it aka JSON.parse()
	jsonResp = "{"
	for i,val := range merchantIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all Merchant")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		fmt.Print("valueAsBytes : ")
		fmt.Println(valueAsBytes)
		jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
		if i < len(merchantIndex)-1 {
			jsonResp = jsonResp + ","
		}
	}
	jsonResp = jsonResp + "}"
	fmt.Println("end getAllMerchants")
	return []byte(jsonResp), nil			//send it onward
}
// ============================================================================================================================
// Delete - remove a merchant from chain
// ============================================================================================================================
func (t *ManageOwner) deleteMerchant(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 'merchantId' as an argument\", \"code\" : \"503\"}"
		err := stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set merchantId
	merchantId := args[0]
	err := stub.DelState(merchantId)													//remove the Merchant from chaincode
	if err != nil {
		errMsg := "{ \"message\" : \"Failed to delete state\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}

	//get the Merchant index
	merchantAsBytes, err := stub.GetState(MerchantIndexStr)
	if err != nil {
		errMsg := "{ \"message\" : \"Failed to get Merchant index\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	var merchantIndex []string
	json.Unmarshal(merchantAsBytes, &merchantIndex)								//un stringify it aka JSON.parse()
	//remove marble from index
	for i,val := range merchantIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + merchantId)
		if val == merchantId{															//find the correct Merchant
			fmt.Println("found Merchant with matching merchantId")
			merchantIndex = append(merchantIndex[:i], merchantIndex[i+1:]...)			//remove it
			for x:= range merchantIndex{											//debug prints...
				fmt.Println(string(x) + " - " + merchantIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(merchantIndex)									//save new index
	err = stub.PutState(MerchantIndexStr, jsonAsBytes)

	tosend := "{ \"merchantID\" : \""+merchantId+"\", \"message\" : \"Merchant deleted succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	fmt.Println("Merchant deleted succcessfully")
	return nil, nil
}
// ============================================================================================================================
// Write - update merchant into chaincode state
// ============================================================================================================================
func (t *ManageOwner) updateMerchant(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("Updating Merchant")
	if len(args) != 8 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 8\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set merchantId
	merchantId := args[0]
	merchantAsBytes, err := stub.GetState(merchantId)									//get the Merchant for the specified merchant from chaincode state
	if err != nil {
		errMsg := "{ \"message\" : \"Failed to get state for " + merchantId + "\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	res := Merchant{}
	json.Unmarshal(merchantAsBytes, &res)
	if res.MerchantID == merchantId{
		fmt.Println("Merchant found with merchantId : " + merchantId)
		//fmt.Println(res);
		res.MerchantName = args[1]
		res.MerchantLogo = args[2]
		res.MerchantIndustry = args[3]
		res.PointsPerDollarSpent = args[4]
		res.MerchantCurrency = args[5]
		res.ExchangeRate = args[6]
		res.MerchantCU_date = args[7]
	}else{
		errMsg := "{ \"message\" : \""+ merchantId+ " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	
	//build the Merchant json string manually
	order := 	`{`+
		`"merchantID": "` + res.MerchantID + `" , `+
		`"merchantName": "` + res.MerchantName + `" , `+
		`"merchantLogo": "` + res.MerchantLogo + `" , `+
		`"merchantIndustry": "` + res.MerchantIndustry + `" , `+ 
		`"pointsPerDollarSpent": "` + res.PointsPerDollarSpent + `" , `+ 
		`"merchantCurrency": "` + res.MerchantCurrency + `" , `+ 
		`"exchangeRate": "` +  res.ExchangeRate + `" , `+ 
		`"merchantCU_date": "` +  res.MerchantCU_date + `" `+ 
		`}`
	err = stub.PutState(merchantId, []byte(order))									//store Merchant with id as key
	if err != nil {
		return nil, err
	}

	tosend := "{ \"merchantId\" : \""+merchantId+"\", \"message\" : \"Merchant details updated succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	fmt.Println("Merchant details updated succcessfully")
	return nil, nil
}
// ============================================================================================================================
// create Merchant - create a new Merchant, store into chaincode state
// ============================================================================================================================
func (t *ManageOwner) createMerchant(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 8 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 8\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Println("start createMerchant")
	merchantID := args[0]
	merchantName := args[1]
	merchantLogo := args[2]
	merchantIndustry := args[3]
	pointsPerDollarSpent := args[4]
	merchantCurrency := args[5]
	exchangeRate := args[6]
	merchantCU_date := args[7]
	merchantAsBytes, err := stub.GetState(merchantID)
	if err != nil {
		return nil, errors.New("Failed to get Merchant merchantID")
	}
	res := Merchant{}
	json.Unmarshal(merchantAsBytes, &res)
	fmt.Print("res: ")
	fmt.Println(res)
	if res.MerchantID == merchantID{
		fmt.Println("This Merchant arleady exists: " + merchantID)
		fmt.Println(res);
		errMsg := "{ \"message\" : \"This Merchant arleady exists\", \"code\" : \"503\"}"
		err := stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil				//all stop a Merchant by this name exists
	}
	//build the Merchant json string manually
	merchant_json := 	`{`+
		`"merchantId": "` + merchantID + `" , `+
		`"merchantName": "` + merchantName + `" , `+
		`"merchantLogo": "` + merchantLogo + `" , `+
		`"merchantIndustry": "` + merchantIndustry + `" , `+ 
		`"pointsPerDollarSpent": "` + pointsPerDollarSpent + `" , `+ 
		`"merchantCurrency": "` + merchantCurrency + `" , `+ 
		`"exchangeRate": "` + exchangeRate + `" , `+ 
		`"merchantCU_date": "` + merchantCU_date + `" `+ 
	`}`
	fmt.Println("merchant_json: " + merchant_json)
	fmt.Print("merchant_json in bytes array: ")
	fmt.Println([]byte(merchant_json))
	err = stub.PutState(merchantID, []byte(merchant_json))									//store Merchant with merchantId as key
	if err != nil {
		return nil, err
	}
	//get the Merchant index
	merchantIndexAsBytes, err := stub.GetState(MerchantIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Merchant index")
	}
	var merchantIndex []string
	json.Unmarshal(merchantIndexAsBytes, &merchantIndex)							//un stringify it aka JSON.parse()
	
	//append
	merchantIndex = append(merchantIndex, merchantID)									//add Merchant merchantID to index list
	fmt.Println("! Merchant index after appending merchantID: ", merchantIndex)
	jsonAsBytes, _ := json.Marshal(merchantIndex)
	fmt.Print("jsonAsBytes: ")
	fmt.Println(jsonAsBytes)
	err = stub.PutState(MerchantIndexStr, jsonAsBytes)						//store name of Merchant
	if err != nil {
		return nil, err
	}

	tosend := "{ \"merchantID\" : \""+merchantID+"\", \"message\" : \"Merchant created succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	fmt.Println("end createMerchant")
	return nil, nil
}