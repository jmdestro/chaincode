// ===================================================================================
// AssetChain v0.1 -  Chaincode definition for assets
// ===================================================================================

package main

// ----- Asset ----- //
type IBMAsset struct {
	ObjectType		   string		`json:"docType"` //field for couchdb
	SerialNumber       string		`json:"serialNumber"`
	Prod               string		`json:"prod"`
	DescriptionProduct string		`json:"descriptionProduct"`
	PwHardware         string		`json:"pwHardware"`
	PwHD               string		`json:"pwHD"`
	PwOS               string		`json:"pwOS"`
	Tickets			   []string		`json:"tickets"`
    OpsUpdate          string  		`json:"opsUpdate"`
	HasWarranty		   string		`json:"hasWarranty"`
}
