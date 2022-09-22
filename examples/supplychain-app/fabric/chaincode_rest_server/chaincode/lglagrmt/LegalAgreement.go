package lglagrmt

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	. "github.com/chaincode/common"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// createLegalAgreement creates an legal agreement in the ledger
func (s *SmartContract) createLegalAgreement(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// Create LegalAgreementRequest struct from input JSON
	argBytes := []byte(args[0])
	var request LegalAgreementRequest
	if err := json.Unmarshal(argBytes, &request); err != nil {
		return shim.Error(fmt.Sprintf("Error unmarshaling LegalAgreementRequest: %s", err))
	}

	// Check if legal agreement state using id as key exists
	testLegalAgreementAsBytes, err := stub.GetState(request.ID)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return 403 if item exists
	if len(testLegalAgreementAsBytes) != 0 {
		return peer.Response{
			Status:  403,
			Message: fmt.Sprintf("Legal Agreement %s already exists", request.ID),
		}
	}

	ContentHash := sha256.Sum256([]byte(request.Content))

	// Call readLatestVersionLegalAgreement to get the latest version
	latestVersionLegalAgreementAsBytes := s.readLatestVersionLegalAgreement(stub, []string{})
	if latestVersionLegalAgreementAsBytes.Status != 200 {
		return shim.Error(fmt.Sprintf("Failed to read latest version legal agreement: %s", latestVersionLegalAgreementAsBytes.Message))
	}

	var latestVersionLegalAgreement LegalAgreement
	err = json.Unmarshal(latestVersionLegalAgreementAsBytes.Payload, &latestVersionLegalAgreement)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Validate that the version is a greater than the previous version
	if latestVersionLegalAgreement.Version >= request.Version {
		return shim.Error(fmt.Sprintf("The version %d is not greater than the latest version %d", request.Version, latestVersionLegalAgreement.Version))
	}

	// Create a new LegalAgreement
	newLegalAgreement := LegalAgreement{
		ID:          request.ID,
		Content:     request.Content,
		ContentHash: fmt.Sprintf("%x", ContentHash),
		Timestamp:   request.Timestamp,
		Version:     request.Version,
	}

	// Marshal legal agreement
	legalAgreementAsBytes, _ := json.Marshal(newLegalAgreement)
	err = stub.PutState(newLegalAgreement.ID, legalAgreementAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	response := map[string]interface{}{
		"createdID": newLegalAgreement.ID,
	}
	bytes, _ := json.Marshal(response)

	s.logger.Infof("Wrote Legal Agreement: %s\n", newLegalAgreement.ID)
	return shim.Success(bytes)
}

// readLegalAgreement returns the legal agreement with the given id
func (s *SmartContract) readLegalAgreement(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// Create ReadLegalAgreementRequest struct from input JSON
	argBytes := []byte(args[0])
	var request ReadLegalAgreementRequest
	if err := json.Unmarshal(argBytes, &request); err != nil {
		return shim.Error(fmt.Sprintf("Error unmarshaling ReadLegalAgreementRequest: %s", err))
	}

	// Get the legal agreement state from the ledger
	legalAgreementAsBytes, err := stub.GetState(request.ID)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return 404 if legal agreement does not exist
	if len(legalAgreementAsBytes) == 0 {
		return peer.Response{
			Status:  404,
			Message: fmt.Sprintf("Legal Agreement %s does not exist", request.ID),
		}
	}

	// Unmarshal legal agreement
	var legalAgreement LegalAgreement
	err = json.Unmarshal(legalAgreementAsBytes, &legalAgreement)
	if err != nil {
		return peer.Response{
			Status:  400,
			Message: fmt.Sprintf("Failed to unmarshal Legal Agreement: %s", err),
		}
	}

	return shim.Success(legalAgreementAsBytes)
}

// readLatestVersionLegalAgreement returns the latest version of the legal agreement
func (s *SmartContract) readLatestVersionLegalAgreement(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Expecting 0")
	}

	// Get iterator for all entries
	iterator, err := stub.GetStateByRange("", "")
	if err != nil {
		shim.Error(fmt.Sprintf("Error getting state iterator: %s", err))
	}
	defer iterator.Close()

	// Get the latest version
	var latestLegalAgreement LegalAgreement
	for iterator.HasNext() {
		// Get the next item
		item, err := iterator.Next()
		if err != nil {
			shim.Error(fmt.Sprintf("Error getting next item: %s", err))
		}

		// Unmarshal item
		var legalAgreement LegalAgreement
		err = json.Unmarshal(item.Value, &legalAgreement)
		if err != nil && err.Error() != "Not a LegalAgreement" {
			shim.Error(fmt.Sprintf("Error unmarshaling item: %s", err))
		}

		// Update latest version
		if latestLegalAgreement.Version < legalAgreement.Version {
			latestLegalAgreement = legalAgreement
		}
	}

	// Marshal latest version
	legalAgreementAsBytes, _ := json.Marshal(latestLegalAgreement)

	return shim.Success(legalAgreementAsBytes)
}
