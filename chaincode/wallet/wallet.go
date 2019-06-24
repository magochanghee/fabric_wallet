// SPDX-License-Identifier: Apache-2.0

/*
  Sample Chaincode based on Demonstrated Scenario

 This code is based on code written by the Hyperledger Fabric community.
  Original code can be found here: https://github.com/hyperledger/fabric-samples/blob/release/chaincode/fabcar/fabcar.go
 */

 package main

 /* Imports  
 * 4 utility libraries for handling bytes, reading and writing JSON, 
 formatting, and string manipulation  
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts  
 */ 
 import (
		 "bytes"
		 "encoding/json"
		 "fmt"
		 "strconv"
		//  "strings"

		 "github.com/hyperledger/fabric/core/chaincode/shim"
		 sc "github.com/hyperledger/fabric/protos/peer"
 )
 
 // Define the Smart Contract structure
type SmartContract struct {
	
}

type Account struct {
	Name string `json:"name"`						//가입자 서명
	Balance int `json:"balance"`					//잔액
	// subscription string `json:"subscription"`		//가입일자
}
/* Define Tuna structure, with 4 properties.  
Structure tags are used by encoding/json library
*/
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
			fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

func (b *SmartContract) Init(stub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (b *SmartContract) Invoke(stub shim.ChaincodeStubInterface) sc.Response{

	function, args := stub.GetFunctionAndParameters()
	if function == "queryBalance"{
		return b.queryBalance(stub, args)
	} else if function == "addUser"{
		return b.addUser(stub, args)
	} else if function == "Remittance"{
		return b.Remittance(stub, args)
	}
	return shim.Error("Invalid invoke funtion name. Expecting \"querybalance\" \"invoke\" \"addUser\"")
}

func (b *SmartContract) addUser(stub shim.ChaincodeStubInterface, args []string) sc.Response{
	if len(args) != 2{
		return shim.Error("Incorrect number of parameters. Expecting 2")
	}

	name := args[0]
	balance, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("2nd argument amount must be a numeric sting \n")
	}
	nameAsBytes, err := stub.GetState(name)
	if err != nil {
		return shim.Error("Faild to get User name " + err.Error())
	} else if nameAsBytes != nil {
		fmt.Println("This User already exists:" + name)
		return shim.Error("This User already exists:" + name)
	}
	account := &Account{Name: args[0],  Balance:balance}
	accountJSONasBytes, _ := json.Marshal(account)
	err = stub.PutState(name, accountJSONasBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failedf to create account :%s", args[0]))
	}
	var buffer bytes.Buffer
	fmt.Println("-end add User Accont -\n")
	buffer.Write(accountJSONasBytes)
	return shim.Success(nil)
}

func (b *SmartContract) queryBalance(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect ID of arguments. Expercting 1")
	}
	balanceAsBytes, _ := stub.GetState(args[0])
	account := Account{}
	json.Unmarshal(balanceAsBytes, &account)

	var buffer bytes.Buffer

	buffer.WriteString("{\"User name\":")
	buffer.WriteString("\"")
	buffer.WriteString(account.Name)
	buffer.WriteString("\"")
	buffer.WriteString(",\n\"amount\":")
	buffer.WriteString(strconv.Itoa(account.Balance))
	buffer.WriteString("}")
	fmt.Printf("Such a balance : \n%s\n", buffer.String())

	if balanceAsBytes == nil {
		return shim.Error("Could not find user data")
	} else {
		return shim.Success(buffer.Bytes())
	}
}

func (t * SmartContract) Remittance(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	// check parameter (send_name, recive_name, amount) 
	if len(args) != 3{
		return shim.Error("Incorrect number of arguments Expecting 3")
	}
	fmt.Println("- Start Remittance -\n")
	SnameBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("faild to get Send_name :" + err.Error())
	} else if SnameBytes == nil {
		return shim.Error("Send name does not exist")
	}
	RnameBytes, err := stub.GetState(args[1])
	if err != nil{
		return shim.Error("faind to get Recive_name :" + err.Error())
	} else if RnameBytes == nil {
		return shim.Error("Recive name does not exist")
	}
	amount, err := strconv.Atoi(args[2])
	if err != nil{
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	Saccount := Account{}
	json.Unmarshal(SnameBytes, &Saccount)
	Samount := Saccount.Balance
	
	Raccount := Account{}
	json.Unmarshal(RnameBytes, &Raccount)
	Ramount  := Raccount.Balance

	if Samount - amount < 0 {
		return shim.Error("Balance to lack")
	} else if amount < 0 {
		return shim.Error("Money is not a minus.")
	}else {
		Samount = Samount - amount
		Ramount = Ramount + amount
		fmt.Println("Samount = %d, Ramount = %d\n", Samount, Ramount)
	}
	Saccount.Balance = Samount
	Raccount.Balance = Ramount

	SnameBytes, err = json.Marshal(Saccount)
	err = stub.PutState(args[0], SnameBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to Remittance : %s", args[0]))
	}
	RnameBytes, err = json.Marshal(Raccount)
	err = stub.PutState(args[1], RnameBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to Remitance : %s", args[1]))
	}
	var buffer bytes.Buffer
	fmt.Println("-end remittance - \n")
	buffer.Write(SnameBytes)
	buffer.Write(RnameBytes)
	return shim.Success(nil)
}