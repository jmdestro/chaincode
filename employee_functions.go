// =================================================
// AssetChain v0.5 - 
// =================================================

package main

import (
	"encoding/json"
	"fmt"
	_ "strconv"
	_ "strings"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ============================================================================================================================
// deleteEmployee() - remove an employee from world state
// ============================================================================================================================

func deleteEmployee(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	fmt.Println("starting deleteEmployee")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// input sanitation
	err := sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	id := args[0]
	
	// get the object
	employee, err := getEmployee(stub, id)
	if err != nil {
		fmt.Println("Failed to find employee by id " + id)
		return shim.Error(err.Error())
	}

	// remove the employee
	err = stub.DelState(employee.UserId) 	//remove the key from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	fmt.Println("- end deleteEmployee")
	return shim.Success(nil)
}

// ============================================================================================================================
// Get Employee - get the employee asset from ledger
// ============================================================================================================================
func getEmployee(stub shim.ChaincodeStubInterface, id string) (Employee, error) {
	
	var employee Employee
	
	employeeAsBytes, err := stub.GetState(id)						//getState retreives a key/value from the ledger
	if err != nil {													//this seems to always succeed, even if key didn't exist
		return employee, errors.New("Failed to get employee - " + id)
	}
	json.Unmarshal(employeeAsBytes, &employee)						//un stringify it aka JSON.parse()

	if len(employee.Email) == 0 {									//test if employee is actually here or just nil
		return employee, errors.New("Employee does not exist - " + id + ", '" + employee.Email + "' '")
	}
	
	return employee, nil  //send it onward
}

// ============================================================================================================================
// saveEmployee - create a new employee and store into chaincode state
// ============================================================================================================================
func saveEmployee(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting init_employee")

	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 8")
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	var employee Employee
	employee.ObjectType = TYPE_EMPLOYEE
	employee.UserId =  args[0]
	employee.Email = args[1]
	employee.Telephone = args[2]
	employee.Cep = args[3]
	employee.Address = args[4]
	employee.City = args[5]
	employee.Uf = args[6]
	employee.OpsUpdate = args[7]

	fmt.Println(employee)	

	//check if employee already exists and updates, or create new
	//_, err = get_employee(stub, employee.UserId)

	//forcedly store a new employee data, even if it already exists
	employeeAsBytes, _ := json.Marshal(employee)	//convert to array of bytes
		err = stub.PutState(employee.UserId, employeeAsBytes)	  //store owner by its Id
		if err != nil {
			fmt.Println("Could not store employee")
			return shim.Error(err.Error())
		}

	fmt.Println("- end init_employee marble")
	return shim.Success(nil)
}

// ============================================================================================================================
// Get history of employee
// ============================================================================================================================
func getEmployeeHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type AuditHistory struct {
		TxId    string   `json:"txId"`
		Value   Employee   `json:"value"`
	}
	var history []AuditHistory;
	var employee Employee

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	employeeId := args[0]
	fmt.Printf("- start getHistoryForTicket: %s\n", employeeId)

	// Get History
	resultsIterator, err := stub.GetHistoryForKey(employeeId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		pointer, err := resultsIterator.Next()
		txID, historicValue := pointer.GetTxId(), pointer.GetValue()
		if err != nil {
			return shim.Error(err.Error())
		}

		var tx AuditHistory
		tx.TxId = txID                             //copy transaction id over
		json.Unmarshal(historicValue, &employee)     //un stringify it aka JSON.parse()
		if historicValue == nil {                  //marble has been deleted
			var emptyEmployee Employee
			tx.Value = emptyEmployee                 //copy nil marble
		} else {
			json.Unmarshal(historicValue, &employee) //un stringify it aka JSON.parse()
			tx.Value = employee                      //copy marble over
		}
		history = append(history, tx)              //add this tx to the list
	}
	fmt.Printf("- getHistoryForEmployee returning:\n%s", history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history)     //convert to array of bytes
	return shim.Success(historyAsBytes)
}
