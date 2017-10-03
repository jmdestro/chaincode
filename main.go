// ===================================================================================
// AssetChain v0.5 -  main
// ===================================================================================

package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode - %s", err)
	}
}

// ============================================================================================================================
// Init - initialize the chaincode - marbles don’t need anything initlization, so let's run a dead simple test instead
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println(" - ready for action")
	return shim.Success(nil)
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println(" ")
	fmt.Println("starting invoke, for - " + function)

	// Handle different functions
	switch function {
		case "init":						//initialize the chaincode state, used as reset
			return t.Init(stub)

		case "read":						//generic read ledger
			return read(stub, args)

		case "read_everything":				//Read Everything
			return read_everything(stub)

		case "readTicket":
			return readTicket(stub, args)

		case "readIBMAsset":
			return readIBMAsset(stub, args)

		case "getTicketsByRange":			//get Ticket by Range
			return getTicketsByRange(stub, args)

		case "getTickets":			//get Ticket by Range
			return getTickets(stub, args)

		case "saveTicket":					//Write Ticket Information
			return saveTicket(stub, args)

		case "saveRCMSTicket":					//Write Ticket Information
			return saveRCMSTicket(stub, args)

		case "getTicketHistory":			//Get History of a ticket by number
			return getTicketHistory(stub, args)

		case "getTicketTimestamp":			//Get the timestamp from the History of a ticket by number
			return getTicketTimestamp(stub, args)

		case "saveIBMAsset":				//Write IBM Asset
			return saveIBMAsset(stub, args)

		case "getTicketsByDate":			// Get all tickets by a time range
			return getTicketsByDate(stub, args)
		default:
			fmt.Println("Received unknown invoke function name - " + function)
			return shim.Error("Received unknown invoke function name - '" + function + "'")
	} //END Switch

} // END func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response

/* PARKING LOT

		case "saveEmployee":				//Write Employee Asset/Ticket Owner
			return saveEmployee(stub, args)

		case "deleteEmployee":				//Delete Employee Asset/Ticket Owner
			return deleteEmployee(stub, args)

		case "getEmployeeHistory": 			//Get history Employee
			return getEmployeeHistory(stub, args)

*/
