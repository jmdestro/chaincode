// =================================================
// AssetChain v0.5 -
// =================================================

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	_ "strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ============================================================================================================================
// Get the timestamp from history of ticket
// ============================================================================================================================
func _getTicketTimestamp(stub shim.ChaincodeStubInterface, ticketId string) string {
	var Timestamp string

	// Get History
	resultsIterator, err := stub.GetHistoryForKey(ticketId)
	if err != nil {
		return ""
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		pointer, err := resultsIterator.Next()
		if err != nil {
			return ""
		}

		Timestamp = time.Unix(pointer.Timestamp.Seconds, int64(pointer.Timestamp.Nanos)).String()
	}

	return Timestamp
}

// ============================================================================================================================
// Get the timestamp from history of ticket
// ============================================================================================================================
func getTicketTimestamp(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var Timestamp string

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ticketId := args[0]
	fmt.Printf("- start getHistoryForTicket: %s\n", ticketId)

	timestamp := _getTicketTimestamp(stub, ticketId)

	fmt.Printf("- getHistoryForTicket returning:\n%s", Timestamp)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(timestamp) //convert to array of bytes
	return shim.Success(historyAsBytes)
}

// ============================================================================================================================
// Get history of ticket
// ============================================================================================================================
func getTicketHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type AuditHistory struct {
		TxId      string `json:"txId"`
		Value     Ticket `json:"value"`
		Timestamp string `json:"timestamp"`
	}
	var history []AuditHistory
	var ticket Ticket

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ticketId := args[0]
	fmt.Printf("- start getHistoryForTicket: %s\n", ticketId)

	// Get History
	resultsIterator, err := stub.GetHistoryForKey(ticketId)
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
		tx.TxId = txID                         //copy transaction id over
		json.Unmarshal(historicValue, &ticket) //un stringify it aka JSON.parse()
		if historicValue == nil {              //marble has been deleted
			var emptyTicket Ticket
			tx.Value = emptyTicket //copy nil marble
		} else {
			json.Unmarshal(historicValue, &ticket) //un stringify it aka JSON.parse()
			if ticket.ObjectType == TYPE_TICKET {
				tx.Value = ticket //copy ticket over
			}
		}
		var ts = time.Unix(pointer.Timestamp.Seconds, int64(pointer.Timestamp.Nanos)).String()
		tx.Timestamp = ts
		history = append(history, tx) //add this tx to the list
	}
	fmt.Printf("- getHistoryForTicket returning:\n%s", history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history) //convert to array of bytes
	return shim.Success(historyAsBytes)
}

// ============================================================================================================================
// Get Ticket - get a ticket asset from ledger
// ============================================================================================================================
func getTicket(stub shim.ChaincodeStubInterface, id string) (Ticket, error) {
	var ticket Ticket

	// input sanitation
	err := sanitize_arguments([]string{id})
	if err != nil {
		return ticket, errors.New(err.Error())
	}

	ticketAsBytes, err := stub.GetState(id) //getState retreives a key/value from the ledger
	if err != nil {                         //this seems to always succeed, even if key didn't exist
		return ticket, errors.New("Failed to find ticket - " + id)
	}
	json.Unmarshal(ticketAsBytes, &ticket) //un stringify it aka JSON.parse()
	if ticket.ObjectType != TYPE_TICKET {
		return ticket, errors.New("id is not a ticket - " + id)
	}

	if ticket.Id != id { //test if ticket is actually here or just nil
		return ticket, errors.New("Ticket does not exist - " + id)
	}

	ticket.Timestamp = _getTicketTimestamp(stub, id)

	return ticket, nil //send it onward
}

// ============================================================================================================================
// Get all of tickets
// ============================================================================================================================
func getTickets(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var tickets []map[string]string

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
		json.Unmarshal(queryValAsBytes, &ticket) //un stringify it aka JSON.parse()

		if ticket.ObjectType == TYPE_TICKET {
			ticket.Timestamp = _getTicketTimestamp(stub, ticket.Id)

			for i := range args {
				if args[i] == "" || args[i] == ticket.StatusTicket {
					// Found!
					ibmasset, _ := getIBMAsset(stub, ticket.SerialNumber)
					out := mergeTicketAsset(ticket, ibmasset)

					tickets = append(tickets, out) //add this Ticket to the list
					break
				}
			}
		}
	}

	fmt.Println("Tickets array - ", tickets)

	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(tickets) //convert to array of bytes
	return shim.Success(everythingAsBytes)
}

// ============================================================================================================================
// Get history of tickets - performs a range query based on the start and end keys provided.
// ============================================================================================================================
func getTicketsByRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	startKey := args[0]
	endKey := args[1]

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		pointer, err := resultsIterator.Next()
		queryResultKey, queryResultValue := pointer.GetKey(), pointer.GetValue()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResultKey)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResultValue))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getTicketsByRange queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// ============================================================================================================================
// Translate Location - method to translate the Location given by the RCMS to a location for Asset Chain
// ============================================================================================================================
func _translateLocation(location string) string {
	var translatedLocation = ""

	switch location {
	case "VENDOR RIO - BUG BUSTERS":
		translatedLocation = "BMA" // Pasteur
	case "VENDOR SAO PAULO - BUG BUSTERS":
		translatedLocation = "BTU" // Tutóia
	case "VENDOR CP - BUG BUSTERS":
		translatedLocation = "BMM" // Hortolândia
	case "VENDOR I/T SITE BIRMANN":
		translatedLocation = "BAW" // Birman
	default:
		translatedLocation = ""
	}

	return translatedLocation
}

// ============================================================================================================================
// Save Ticket - method to expose to assetchain and create a new ticket or modify it, store into chaincode state
// ============================================================================================================================
func saveRCMSTicket(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting saveTicket")

	if len(args) != 14 {
		return shim.Error("Incorrect number of arguments. Expecting 14")
	}

	//build ticket object
	var ticket Ticket
	ticket.ObjectType = TYPE_TICKET
	ticket.Id = args[0]
	ticket.OpenDate = args[1]
	ticket.StatusTicket = args[2]
	ticket.Queue = args[3]
	ticket.DescriptionRequest = args[4]
	ticket.UserRequest = args[5]
	ticket.Telephone = args[6]
	ticket.Cep = args[7]
	ticket.Address = args[8]
	ticket.City = args[9]
	ticket.SerialNumber = args[10]
	ticket.OpsUpdate = USER_RCMS
	ticket.OpsAssignTicket = ""
	ticket.Email = ""
	ticket.Uf = ""
	ticket.Logbook = ""
	ticket.TokenIn = ""
	ticket.TokenOut = ""
	ticket.IbmLocation = _translateLocation(args[13]) // Gets the Location from RCMS and Translate to a location to be used on Asset chain
	ticket.RepairType = []string{""}

	//build ibmasset array
	var ibmasset = IBMAsset{
		SerialNumber:       ticket.SerialNumber,
		Prod:               args[11],
		DescriptionProduct: args[12],
		PwHardware:         "",
		PwHD:               "",
		PwOS:               "",
		Tickets:            []string{ticket.Id},
		OpsUpdate:          ticket.OpsUpdate,
		HasWarranty:        "",
		Status:        "Active"}

	_saveIBMAsset(stub, ibmasset)

	//check if ticket exists
	existingTicket, err := getTicket(stub, ticket.Id)
	if err == nil {
		// If ticket already exists and data came from RCMS, keep relevant information from blockchain
		if ticket.OpsUpdate == USER_RCMS {
			if existingTicket.StatusTicket == STATUS_RECEBIDO || existingTicket.StatusTicket == STATUS_EXECUCAO || existingTicket.StatusTicket == STATUS_CONCLUIDO || existingTicket.StatusTicket == STATUS_ENCERRADO {
				ticket.StatusTicket = existingTicket.StatusTicket
			}

			ticket.DescriptionRequest = existingTicket.DescriptionRequest
			ticket.UserRequest = existingTicket.UserRequest
			ticket.Email = existingTicket.Email
			ticket.Telephone = existingTicket.Telephone
			ticket.Cep = existingTicket.Cep
			ticket.Address = existingTicket.Address
			ticket.City = existingTicket.City
			ticket.Uf = existingTicket.Uf
			ticket.OpsUpdate = existingTicket.OpsUpdate
			ticket.Logbook = existingTicket.Logbook
			ticket.TokenIn = existingTicket.TokenIn
			ticket.TokenOut = existingTicket.TokenOut
			ticket.IbmLocation = existingTicket.IbmLocation
			ticket.RepairType = existingTicket.RepairType
		}
	}

	fmt.Println(ticket)

	var ticketMessage, ticketErr = validateTicket(ticket)
	if ticketErr == nil {
		//forcedly store a new ticket data, even if it already exists
		ticketAsBytes, _ := json.Marshal(ticket)      //convert to array of bytes
		err = stub.PutState(ticket.Id, ticketAsBytes) //store owner by its Id
		if err != nil {
			fmt.Println("Could not store ticket")
			return shim.Error(err.Error())
		}
	} else {
		// Ticket has problem based on rules
		fmt.Println(ticketMessage)
		return shim.Error(ticketMessage)
	}

	fmt.Println("end saveTicket")
	return shim.Success(nil)
}

// ============================================================================================================================
// Save Ticket - method to expose to assetchain and create a new ticket or modify it, store into chaincode state
// ============================================================================================================================
func saveTicket(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting saveTicket")

	if len(args) != 27 {
		return shim.Error("Incorrect number of arguments. Expecting 27")
	}

	//build ticket object
	var ticket Ticket
	ticket.ObjectType = TYPE_TICKET
	ticket.Id = args[0]
	ticket.OpenDate = args[1]
	ticket.StatusTicket = args[2]
	ticket.Queue = args[3]
	ticket.OpsAssignTicket = args[4]
	ticket.DescriptionRequest = args[5]
	ticket.UserRequest = args[6]
	ticket.SerialNumber = args[13]
	ticket.Email = args[7]
	ticket.Telephone = args[8]
	ticket.Cep = args[9]
	ticket.Address = args[10]
	ticket.City = args[11]
	ticket.Uf = args[12]
	ticket.OpsUpdate = args[19]
	ticket.Logbook = args[20]
	ticket.TokenIn = args[21]
	ticket.TokenOut = args[22]
	ticket.IbmLocation = args[23]

	repairType := strings.Split(args[25], ",")
	if len(repairType) > 0 {
		ticket.RepairType = repairType
	} else {
		ticket.RepairType = []string{""}
	}

	//build ibmasset array
	var ibmasset = IBMAsset{
		SerialNumber:       ticket.SerialNumber,
		Prod:               args[14],
		DescriptionProduct: args[15],
		PwHardware:         args[16],
		PwHD:               args[17],
		PwOS:               args[18],
		Tickets:            []string{ticket.Id},
		OpsUpdate:          ticket.OpsUpdate,
		HasWarranty:        args[24],
		Status:        		args[26]}

	_saveIBMAsset(stub, ibmasset)

	fmt.Println(ticket)

	var ticketMessage, ticketErr = validateTicket(ticket)
	if ticketErr == nil {
		//forcedly store a new ticket data, even if it already exists
		ticketAsBytes, _ := json.Marshal(ticket)      //convert to array of bytes
		err = stub.PutState(ticket.Id, ticketAsBytes) //store owner by its Id
		if err != nil {
			fmt.Println("Could not store ticket")
			return shim.Error(err.Error())
		}
	} else {
		// Ticket has problem based on rules
		fmt.Println(ticketMessage)
		return shim.Error(ticketMessage)
	}

	fmt.Println("end saveTicket")
	return shim.Success(nil)
}

// ============================================================================================================================
// deleteTicket() - remove a ticket from state and from ticket index
// ============================================================================================================================
func deleteTicket(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("starting delete_ticket")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	id := args[0]

	// get the object
	ticket, err := getTicket(stub, id)
	if err != nil {
		fmt.Println("Failed to find ticket by id " + id)
		return shim.Error(err.Error())
	}

	// remove the ticket
	err = stub.DelState(ticket.Id) //remove the key from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	fmt.Println("- end delete_ticket")
	return shim.Success(nil)
}

// ============================================================================================================================
// Set Assignee on Ticket
// ============================================================================================================================
func set_assignee(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting set_assignee")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	var ticket_id = args[0]
	var assignee = args[1]
	fmt.Println(ticket_id + "->" + assignee)

	// check if user already exists
	//employee, err := get_employee(stub, assignee)
	if err != nil {
		return shim.Error("This employee does not exist - " + assignee)
	}

	// get ticket's current state
	ticketAsBytes, err := stub.GetState(ticket_id)
	if err != nil {
		return shim.Error("Failed to get ticket")
	}
	res := Ticket{}
	json.Unmarshal(ticketAsBytes, &res) //un stringify it aka JSON.parse()

	// set assignee
	//res.Assignee.Employee_sn = employee.Employee_sn               //change the assignee
	jsonAsBytes, _ := json.Marshal(res)       //convert to array of bytes
	err = stub.PutState(args[0], jsonAsBytes) //rewrite the ticket with id as key
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end set assignee")
	return shim.Success(nil)
}

// ============================================================================================================================
// ReadTicket - read a ticket from ledger with asset information
// ============================================================================================================================
func readTicket(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, jsonResp string
	var err error
	fmt.Println("starting read")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting key of the var to query")
	}

	key = args[0]
	ticket, err := getTicket(stub, key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get asset for " + key + "\"}"
		return shim.Error(jsonResp)
	}

	//get asset from ledger
	ibmasset, err := getIBMAsset(stub, ticket.SerialNumber)
	if err != nil {
		fmt.Println("Failed to find asset by id " + ticket.SerialNumber)
		return shim.Error(err.Error())
	}

	//encode the output
	valAsbytes, err := json.Marshal(mergeTicketAsset(ticket, ibmasset))
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get asset for " + key + "\"}"
		return shim.Error(jsonResp)
	}

	fmt.Println("- end read")
	return shim.Success(valAsbytes) //send it onward
}

func mergeTicketAsset(ticket Ticket, ibmasset IBMAsset) map[string]string {
	// merge IBMAsset and Tickets
	var out map[string]string

	tempAsset, _ := json.Marshal(ibmasset)
	json.Unmarshal(tempAsset, &out)
	//remove the field Tickets from asset
	delete(out, "tickets")

	var repairType = strings.Join(ticket.RepairType, ",")
	tempTicket, _ := json.Marshal(ticket)
	json.Unmarshal(tempTicket, &out)
	delete(out, "repairType")
	out["repairType"] = repairType

	//encode the output
	return out
}

// ============================================================================================================================
// Get all tickets by Date
// ============================================================================================================================
func getTicketsByDate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var tickets []map[string]string

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	startDate := args[0]
	endDate := args[1]

	start, _ := time.Parse("2006-01-02", startDate)
    end, _ := time.Parse("2006-01-02", endDate)

	start = start.AddDate(0, 0, -1)
	end = end.AddDate(0, 0, 1)

	// fmt.Printf("start: %d-%02d-%02d\n", start.Year(), start.Month(), start.Day())
	// fmt.Printf("end: %d-%02d-%02d\n", end.Year(), end.Month(), end.Day())

	// ---- Get All Tickets ---- //
	resultsIterator, err := stub.GetStateByRange("", "z999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		pointer, err := resultsIterator.Next()
		_, queryValAsBytes := pointer.GetKey(), pointer.GetValue()
	    
		if err != nil {
			return shim.Error(err.Error())
		}

		// fmt.Println("on ticket id - ", queryKeyAsStr)
		var ticket Ticket
		json.Unmarshal(queryValAsBytes, &ticket) //un stringify it aka JSON.parse()

		if ticket.ObjectType == TYPE_TICKET {

			openDate := strings.Split(ticket.OpenDate, " ")[0]
			filter := "2006-01-02"
			if (strings.Contains(openDate, "/")) {
				filter = "2006/01/02"
			}

			in, _ := time.Parse(filter, openDate)

			if inTimeSpan(start, end, in) {
				ibmasset, _ := getIBMAsset(stub, ticket.SerialNumber)
				out := mergeTicketAsset(ticket, ibmasset)

				tickets = append(tickets, out) //add this Ticket to the list
			}
		}
	}

	// fmt.Println("Tickets array - ", tickets)

	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(tickets) //convert to array of bytes
	return shim.Success(everythingAsBytes)
}

func inTimeSpan(start, end, check time.Time) bool {
    return check.After(start) && check.Before(end)
}
