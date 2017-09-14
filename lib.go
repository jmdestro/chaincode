// =================================================
// AssetChain v0.1 - get ticket, employee and asset
// =================================================

package main

import (
	"errors"
	"strconv"
)

const (
	USER_RCMS			= "RCMS"
	TYPE_TICKET			= "ticket"
	TYPE_ASSET			= "ibm_asset"
	TYPE_EMPLOYEE		= "employee"
	STATUS_RECEBIDO		= "RECEB"
	STATUS_EXECUCAO		= "EXEC"
	STATUS_CONCLUIDO	= "CONCL"
	STATUS_ENCERRADO	= "CLO"
)

// ========================================================
// Input Sanitation - dumb input checking, look for empty strings
// ========================================================
func sanitize_arguments(strs []string) error{
	for i, val:= range strs {
		if len(val) <= 0 {
			return errors.New("Argument " + strconv.Itoa(i) + " must be a non-empty string")
		}
		if len(val) > 32 {
			return errors.New("Argument " + strconv.Itoa(i) + " must be <= 32 characters")
		}
	}
	return nil
}

// ========================================================
// Unique Element - function to add an unique element to slice
// ========================================================
func uniqueElem(sliceMain []string, elem string) []string{
  for _, elemSliceMain := range sliceMain {
        if elemSliceMain == elem {
            return sliceMain
        }
    }
    return append(sliceMain, elem)
}

