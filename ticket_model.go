// ===================================================================================
// AssetChain v0.5 -  Chaincode definition for tickets
// ===================================================================================

package main

// ----- Tickets ----- //
type Ticket struct {
	ObjectType         string  `json:"docType"`     //field for couchdb
	Id                 string  `json:"id"`
	OpenDate           string  `json:"openDate"`
	StatusTicket       string  `json:"statusTicket"`
	Queue              string  `json:"queue"`
	OpsAssignTicket    string  `json:"opsAssignTicket"`
	DescriptionRequest string  `json:"descriptionRequest"`
	UserRequest        string  `json:"userRequest"`
	SerialNumber       string  `json:"serialNumber"`	
	Email              string  `json:"email"`
	Telephone          string  `json:"telephone"`
	Cep                string  `json:"cep"`
	Address            string  `json:"address"`
	City               string  `json:"city"`
	Uf                 string  `json:"uf"`
    OpsUpdate          string  `json:"opsUpdate"`
    Logbook            string  `json:"logbook"`
    TokenIn            string  `json:"tokenIn"`
    TokenOut           string  `json:"tokenOut"`
	IbmLocation		   string  `json:"ibmLocation"`
	RepairType		   []string `json:"repairType"`
	Timestamp		   string	`json:"timestamp",omitempty`
}
