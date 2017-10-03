// =================================================
// AssetChain v0.5 - Asset SmartContract
// =================================================

package main

import (
	"fmt"
)

// ============================================================================================================================
// ValidateAssetStatus - validate the status change before saving on the chain
// ============================================================================================================================
func validateAssetStatus(ibmasset IBMAsset, newStatus string) (bool, error) {
	switch newStatus {
		case "BDC":	
			if ibmasset.Status == "Active" || ibmasset.Status == ""{
				return true, nil
			} else
			{
				return false, nil
			}
		case "Active":
			if ibmasset.Status == "BDC" || ibmasset.Status == ""{
				return true, nil
			} else
			{
				return false, nil
			}
		default:
			fmt.Println("Received unknown status  " + newStatus)
			return false, fmt.Errorf("Received unknown status - '" + newStatus + "'")
	} //END Switch

	return false, nil
}
