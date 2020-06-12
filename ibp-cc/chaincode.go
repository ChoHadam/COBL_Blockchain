/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Chaincode is the definition of the chaincode structure.
type Chaincode struct {
}

//TestModel is the definition of the chaincode structure.
type donation struct{
	OrderId string
	ContributerId string
	InstitutionId string
	AFactoryId string
	BFactoryId string
	Status string
	DonateDate string
	transferDate string
	CompleteDate string //공장에서 물품 오더를 성공적으로 받음
	DonatePrice int
	GoalPrice int
	CurrentAmount int	
}	

// Init is called when the chaincode is instantiated by the blockchain network.
func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) sc.Response {
	return setupChainOrder(stub)
}

// Invoke is called as a result of an application request to run the chaincode.
func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "createDonation" {
		return t.createDonation(stub, args)
	} else if function == "checkDonation" { //총 모금액 체크

		return t.checkDonation(stub, args)
	} else if function == "initiateTransfer" { //송금 초기화

		return t.initiateTransfer(stub, args)
	} else if function == "TransferToFactory" { //

		return t.TransferToFactory(stub, args)
	} else if function == "completeOrder" { 

		return t.completeOrder(stub, args)
	} else if function == "query" {

		return t.query(stub, args)
	}
	return shim,Error("Invalid Chaincode function name")
}

//
func setupChainOrder(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	_, args := stub.GetFunctionAndParameters()
	orderId := args[0]
	contributerId := args[1]
	currentAmount := args[2]
	donatePrice, _ := strconv.Atoi(args[3])

	chainCode := donation{
		OrderId :orderId
		ContributerId : contributerId
		CurrentAmount : 0
		DonatePrice : donatePrice
		Status : "order initiated"}
	donationBytes, _ := json.Marshal(chainCode)
	stub.PutState(chainCode.OrderId, donationBytes)

	return shim.Success(nil)
}

//
func (cc *Chaincode) createDonation(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	orderId := args[0]
	donationBytes, _ := stub.GetState(orderId)
	dn := donation{}
	json.Unmarshal(donationBytes, &dn)

	if dn.Status == "order initiated" {
		currentts := time.Now()
		dn.DonateDate = currentts.Format("2020-06-13 10:04:05")
		dn.Status = "donation created"
	} else {
		fmt.Printf("Donation not initiated yet")
	}

	donationBytes, _ := json.Marshal(dn)
	stub.PutState(orderId, donationBytes)

	return shim.Success(nil)
}

func (cc *Chaincode) addDonation(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	orderId := args[0]
	donationBytes, err := stub.GetState(orderId)
	goalPrice := args[1]
	currentAmount := args[2]
	donatePrice, _ := strconv.Atoi(args[3])

	dn := donation{}
	err = json.Unmarshal(donationBytes, &dn)
	if err != nil {
		return shim.Error(err.Error())
	}

	if dn.Status == "donation added" {
		dn.InstitutionId = "InstitutionId_1"
		currentts := time.Now()
		dn.GoalPrice = goalPrice
		dn.CurrentAmount = currentAmount +donatePrice
		dn.Status = "addDonation Process"
	} else {
		dn.Status = "Error"
		fmt.Printf("Create Donation not initiated yet")
	}

	donationBytes0, _ := json.Marshal(dn)
	err = stub.PutState(orderId, donationBytes0)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (cc *Chaincode) initiateTransfer(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	orderId := args[0]
	donationBytes, err := stub.GetState(orderId)
	dn := donation{}
	err = json.Unmarshal(donationBytes, &dn)
	if err != nil {
		return shim.Error(err.Error())
	}
	if CurrentAmount < GoalPrice {
		fmt.Printf("Wholesaler not initiated yet")
	}
	else {
		dn.Status == "wholesaler distribute" {
			dn.Status = "initiated transfer"
		} else {
			fmt.Printf("Create Donation not initiated yet")
		}
	}	

	donationBytes0, _ := json.Marshal(dn)
	err = stub.PutState(orderId, donationBytes0)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (cc *Chaincode) TransferToFactory(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	orderId := args[0]
	donationBytes, err := stub.GetState(orderId)
	dn := donation{}
	err = json.Unmarshal(donationBytes, &dn)
	if err != nil {
		return shim.Error(err.Error())
	}

	if dn.Status == "initiated transfer" {
		dn.AFactoryId = "Factory_A"
		dn.BFactoryId = "Factory_B"
		currentts := time.Now()
		dn.transferDate = currentts.Format("2020-06-13 16:04:05")
		dn.Status = "Transfer started"

	} else {
		fmt.Printf("Transfer not initiated yet")
	}

	donationBytes0, _ := json.Marshal(dn)
	err = stub.PutState(orderId, donationBytes0)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (cc *Chaincode) completeOrder(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	orderId := args[0]
	donationBytes, err := stub.GetState(orderId)
	dn := donation{}
	err = json.Unmarshal(donationBytes, &dn)
	if err != nil {
		return shim.Error(err.Error())
	}

	if dn.Status == "Transfer started" {
		currentts := time.Now()
		dn.completeDate = currentts.Format("2006-01-02 15:04:05")
		dn.Status = "Factory received order"
	} else {
		fmt.Printf("Transfer to Factory not started yet")
	}

	DONATIONBytes0, _ := json.Marshal(DN)
	err = stub.PutState(orderId, donationBytes0)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (cc *ChainCode) query(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	var ENIITY string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expected ENIITY Name")
	}

	ENIITY = args[0]
	Avalbytes, err := stub.GetState(ENIITY)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + ENIITY + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil order for " + ENIITY + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(Avalbytes)
}

func main() {

	err := shim.Start(new(ChainCode))
	if err != nil {
		fmt.Printf("Error creating new Chain Code: %s", err)
	}
}