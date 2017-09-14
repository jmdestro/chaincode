// ===================================================================================
// AssetChain v0.1 -  Chaincode definition for employees
// ===================================================================================

package main

// ----- Employee ----- //
type Employee struct {
	ObjectType	       string `json:"docType"`     //field for couchdb
	UserId	           string `json:"userId"`
	Email              string `json:"email"`
	Telephone          string `json:"telephone"`
	Cep                string `json:"cep"`
	Address            string `json:"address"`
	City               string `json:"city"`
	Uf                 string `json:"uf"`
    OpsUpdate          string `json:"opsUpdate"`
}
