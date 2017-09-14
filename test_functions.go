/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"fmt"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func queryTest(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	testAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(testAsBytes)
}

func createTest(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	var test = Test{Status: args[1]}

	testAsBytes, _ := json.Marshal(test)
	APIstub.PutState(args[0], testAsBytes)

	return shim.Success(nil)
}

// ============================================================================================================================
// Get everything we need (tickets + employee + asset)
//
// Inputs - none
//
// Returns: json array with tickets, employees and assets
// ============================================================================================================================
func read_everything(stub shim.ChaincodeStubInterface) sc.Response {
	type Everything struct {
		Tickets   []Ticket   `json:"tickets"`
	//	Employees  []Employee  `json:"employee"`
		Assets  []IBMAsset  `json:"ibmasset"`
	}
	var everything Everything

	// ---- Get All Tickets ---- //
	resultsIterator, err := stub.GetStateByRange("", "z999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
	
	for resultsIterator.HasNext() {
		pointer, err := resultsIterator.Next()
		queryKeyAsStr, queryValAsBytes := pointer.GetKey(), pointer.GetValue()
		fmt.Println("queryKeyAsStr =" + queryKeyAsStr)
		if err != nil {
			return shim.Error(err.Error())
		}

		fmt.Println("on ticket id - ", queryKeyAsStr)
		var ticket Ticket
		json.Unmarshal(queryValAsBytes, &ticket)                  //un stringify it aka JSON.parse()
		if ticket.ObjectType == "ticket" {
			ticket.Timestamp = _getTicketTimestamp(stub, ticket.Id)
			everything.Tickets = append(everything.Tickets, ticket)   //add this Ticket to the list
		}
	}
	fmt.Println("Tickets array - ", everything.Tickets)

	// ---- Get All Employees ---- //
	/*employeeIterator, err := stub.GetStateByRange("0", "9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer employeeIterator.Close()

	for employeeIterator.HasNext() {
		 pointer, err := employeeIterator.Next()
		 queryKeyAsStr, queryValAsBytes := pointer.GetKey(), pointer.GetValue()
		 fmt.Println("queryKeyAsStr =" + queryKeyAsStr)
		if err != nil {
			return shim.Error(err.Error())
		}
		
		fmt.Println("on employee id - ", queryKeyAsStr)
		var employee Employee
		json.Unmarshal(queryValAsBytes, &employee)                  //un stringify it aka JSON.parse()
		if employee.ObjectType == "employee" {
			everything.Employees = append(everything.Employees, employee)     //add this Employee to the list
		}
	}
	fmt.Println("Employees array - ", everything.Employees)
*/
	// ---- Get All Assets ---- //
	assetIterator, err := stub.GetStateByRange("0", "9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer assetIterator.Close()

	for assetIterator.HasNext() {
		pointer, err := assetIterator.Next()
		queryKeyAsStr, queryValAsBytes := pointer.GetKey(), pointer.GetValue()
		if err != nil {
			return shim.Error(err.Error())
		}
		
		fmt.Println("on asset id - ", queryKeyAsStr)
		var ibmasset IBMAsset
		json.Unmarshal(queryValAsBytes, &ibmasset)                  //un stringify it aka JSON.parse()
		if ibmasset.ObjectType == "ibm_asset" {
			everything.Assets = append(everything.Assets, ibmasset)     //add this asset to the list
		}
	}
	fmt.Println("Assets array - ", everything.Assets)

	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(everything)             //convert to array of bytes
	return shim.Success(everythingAsBytes)
}

// ============================================================================================================================
// Read - read a generic variable from ledger
// ============================================================================================================================
func read(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	var key, jsonResp string
	var err error
	fmt.Println("starting read")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting key of the var to query")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)           //get the var from ledger
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(jsonResp)
	}

	fmt.Println("- end read")
	return shim.Success(valAsbytes)                  //send it onward
}
