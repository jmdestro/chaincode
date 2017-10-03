// =================================================
// AssetChain v0.5 - write objects to ledger
// =================================================

package main

import (
	"encoding/json"
	"fmt"
	_ "strconv"
	_ "strings"
	"strings"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ============================================================================================================================
// Get IBMAsset - get the IBMAsset asset from ledger
// ============================================================================================================================
func getIBMAsset(stub shim.ChaincodeStubInterface, id string) (IBMAsset, error) {
	var ibmasset IBMAsset

	// input sanitation
	err := sanitize_arguments([]string {id})
	if err != nil {
		return ibmasset, errors.New(err.Error())
	}

	ibmassetAsBytes, err := stub.GetState(id)                     //getState retreives a key/value from the ledger
	if err != nil {                                            //this seems to always succeed, even if key didn't exist
		return ibmasset, errors.New("Failed to get Asset - " + id)
	}
	json.Unmarshal(ibmassetAsBytes, &ibmasset)                       //un stringify it aka JSON.parse()
	if ibmasset.ObjectType != TYPE_ASSET {
		return ibmasset, errors.New("Asset retrieved is not a ibm_asset - " + id)
	}

	if len(ibmasset.SerialNumber) == 0 {                              //test if asset is actually here or just nil
		return ibmasset, errors.New("Asset does not exist - " + id)
	}

	return ibmasset, nil                  //send it onward
}

// ============================================================================================================================
// Save Asset - create a new IBMAsset and store into chaincode state
// ============================================================================================================================
func _saveIBMAsset(stub shim.ChaincodeStubInterface, ibmasset IBMAsset) pb.Response {
	var err error

	ibmasset.ObjectType			= TYPE_ASSET

	//check if ibmasset exists to append ticket
	//var arrayTickets []string
	existingIBMAsset, err := getIBMAsset(stub, ibmasset.SerialNumber)


	
	if err == nil {
		
		if ibmasset.Status == "" {
			if existingIBMAsset.Status != "" {
				ibmasset.Status = existingIBMAsset.Status
			} else {
				ibmasset.Status = "Active"
			}
		}

		statusValidation, err := validateAssetStatus(existingIBMAsset, ibmasset.Status)
		
		if (strings.Compare(existingIBMAsset.Status, ibmasset.Status) == 0) || (statusValidation && err == nil) {
		
			for i := range ibmasset.Tickets {
				ibmasset.Tickets = uniqueElem(existingIBMAsset.Tickets, ibmasset.Tickets[i])
				
			}

			// If asset already exists and data came from RCMS, keep relevant information from blockchain
			if ibmasset.OpsUpdate == USER_RCMS {
				ibmasset.Prod				= existingIBMAsset.Prod
				ibmasset.DescriptionProduct	= existingIBMAsset.DescriptionProduct
				ibmasset.PwHardware			= existingIBMAsset.PwHardware
				ibmasset.PwHD				= existingIBMAsset.PwHD
				ibmasset.PwOS				= existingIBMAsset.PwOS
				ibmasset.OpsUpdate			= existingIBMAsset.OpsUpdate
				ibmasset.HasWarranty		= existingIBMAsset.HasWarranty
				ibmasset.Status			    = existingIBMAsset.Status
			}
		} else {
			return shim.Error("Invalid status change from " + existingIBMAsset.Status + " to " + ibmasset.Status)
		}
	}

	//store asset

	ibmassetAsBytes, _ := json.Marshal(ibmasset)	//convert to array of bytes
	err = stub.PutState(ibmasset.SerialNumber, ibmassetAsBytes)	  //store asset by its serial number
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// ============================================================================================================================
// Save Asset - create a new IBMAsset and store into chaincode state
// ============================================================================================================================
func saveIBMAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("starting init_ibmasset")

	if len(args) != 10 {
		return shim.Error("Incorrect number of arguments. Expecting 10")
	}

	var ibmasset IBMAsset
	ibmasset.SerialNumber		= args[0]
	ibmasset.Prod				= args[1]
	ibmasset.DescriptionProduct	= args[2]
	ibmasset.PwHardware			= args[3]
	ibmasset.PwHD				= args[4]
	ibmasset.PwOS				= args[5]
	ibmasset.Tickets			= append(ibmasset.Tickets, args[6])
	ibmasset.OpsUpdate			= args[7]
	ibmasset.HasWarranty	    = args[8]
	if args[9] == "" {
		ibmasset.Status = "Active"
	} else {
		ibmasset.Status = args[9]
	}
	
	fmt.Println(ibmasset)

	return _saveIBMAsset(stub, ibmasset)
}

// ============================================================================================================================
// deleteIBMAsset() - remove an IBMAsset from state and from IBMAsset index
// ============================================================================================================================
func deleteIBMAsset(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	fmt.Println("starting delete_ibmasset")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// input sanitation
	err := sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	id := args[0]

	// get the object
	ibmasset, err := getIBMAsset(stub, id)
	if err != nil {
		fmt.Println("Failed to find asset by id " + id)
		return shim.Error(err.Error())
	}

	// remove the Asset
	err = stub.DelState(ibmasset.SerialNumber)   	 //remove the key from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	fmt.Println("- end delete_ibmasset")
	return shim.Success(nil)
}

// ============================================================================================================================
// ReadIBMAsset - read an ibmasset from ledger and return an array with ticket information
// ============================================================================================================================
func readIBMAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var id, jsonResp string
	var err error

	type AssetDetails struct {
		IbmAsset		IBMAsset	`json:"ibmAsset"`
		TicketsDetails	[]Ticket	`json:"ticketsDetails"`
	}

	var assetDetails AssetDetails

	fmt.Println("starting read")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	id = args[0]

	ibmasset, err := getIBMAsset(stub, id)
	assetDetails.IbmAsset = ibmasset
	if err != nil {
		fmt.Println("Failed to find asset by id " + id)
		return shim.Error(err.Error())
	}

	//get tickets from ledger and add to assetDetails.TicketsDetails
	for _, eachTicket := range assetDetails.IbmAsset.Tickets {

		ticket, err := getTicket(stub, eachTicket)
		if err == nil {
			assetDetails.TicketsDetails = append(assetDetails.TicketsDetails, ticket)
		}
	}

	//jsonify assetDetails to send it onward
	valAsbytes, err := json.Marshal(assetDetails)
	if err != nil {
		jsonResp = "{\"Error\":\" Failed to convert ibmasset " + id + "\"}"
		return shim.Error(jsonResp)
	}

	fmt.Println("- end read")
	return shim.Success(valAsbytes)                  //send it onward
}
