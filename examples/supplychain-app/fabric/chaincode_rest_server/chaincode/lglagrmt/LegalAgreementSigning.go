package lglagrmt

import (
	"encoding/json"
	"fmt"

	. "github.com/chaincode/common"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// createLegalAgreementSigning creates a legal agreement signing in the ledger
func (s *SmartContract) createLegalAgreementSigning(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// Create LegalAgreementSigningRequest struct from input JSON
	argBytes := []byte(args[0])
	var request LegalAgreementSigningRequest
	if err := json.Unmarshal(argBytes, &request); err != nil {
		return shim.Error(fmt.Sprintf("Error unmarshaling LegalAgreementSigningRequest: %s", err))
	}

	// Check if legal agreement signing state using id as key exists
	testLegalAgreementSigningAsBytes, err := stub.GetState(request.ID)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return 404 if result's empty
	if len(testLegalAgreementSigningAsBytes) != 0 {
		return peer.Response{
			Status:  403,
			Message: fmt.Sprintf("Legal Agreement Signing %s already exists", request.ID),
		}
	}

	// Call readLegalAgreement to get the latest version
	readLegalAgreementRequest := &ReadLegalAgreementRequest{ID: request.LegalAgreementID}
	readLegalAgreementRequestAsBytes, err := json.Marshal(readLegalAgreementRequest)
	if err != nil {
		return shim.Error(fmt.Sprintf("Error marshaling ReadLegalAgreementRequest: %s", err))
	}
	legalAgreementAsBytes := s.readLegalAgreement(stub, []string{string(readLegalAgreementRequestAsBytes)})
	if legalAgreementAsBytes.Status != 200 {
		return shim.Error(fmt.Sprintf("Failed to read legal agreement: %s", legalAgreementAsBytes.Message))
	}

	// Check if content hash is equal to the latest version
	var legalAgreement LegalAgreement
	err = json.Unmarshal(legalAgreementAsBytes.Payload, &legalAgreement)
	if err != nil {
		return shim.Error(err.Error())
	}

	if legalAgreement.ContentHash != request.LegalAgreementContentHash {
		return shim.Error(fmt.Sprintf("Content hash does not match latest version of legal agreement"))
	}

	// Create a new LegalAgreementSigning
	newLegalAgreementSigning := LegalAgreementSigning{
		ID:                        request.ID,
		UserID:                    request.UserID,
		LegalAgreementID:          request.LegalAgreementID,
		LegalAgreementContentHash: request.LegalAgreementContentHash,
		Accepted:                  request.Accepted,
		Timestamp:                 request.Timestamp,
	}

	// Marshal legal agreement signing
	legalAgreementSigningAsBytes, _ := json.Marshal(newLegalAgreementSigning)
	err = stub.PutState(newLegalAgreementSigning.ID, legalAgreementSigningAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	response := map[string]interface{}{
		"createdID": newLegalAgreementSigning.ID,
	}
	bytes, _ := json.Marshal(response)

	s.logger.Infof("Wrote Legal Agreement Signing: %s\n", newLegalAgreementSigning.ID)
	return shim.Success(bytes)
}

// readLegalAgreementSigning returns the legal agreement signing with the given id
func (s *SmartContract) readLegalAgreementSigning(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// Create ReadLegalAgreementSigningRequest struct from input JSON
	argBytes := []byte(args[0])
	var request ReadLegalAgreementSigningRequest
	if err := json.Unmarshal(argBytes, &request); err != nil {
		return shim.Error(fmt.Sprintf("Error unmarshaling ReadLegalAgreementSigningRequest: %s", err))
	}

	// Get the legal agreement signing state from the ledger
	legalAgreementSigningAsBytes, err := stub.GetState(request.ID)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return 404 if result's empty
	if len(legalAgreementSigningAsBytes) == 0 {
		return peer.Response{
			Status:  404,
			Message: fmt.Sprintf("Legal Agreement Signing %s does not exist", request.ID),
		}
	}

	// Unmarshal legal agreement signing
	var legalAgreementSigning LegalAgreementSigning
	err = json.Unmarshal(legalAgreementSigningAsBytes, &legalAgreementSigning)
	if err != nil {
		return peer.Response{
			Status:  400,
			Message: fmt.Sprintf("Error unmarshalling Legal Agreement Signing: %s", err.Error()),
		}
	}

	return shim.Success(legalAgreementSigningAsBytes)
}

// readLatestLegalAgreementSigningByUserID returns the latest legal agreement signing by user id
func (s *SmartContract) readLatestLegalAgreementSigningByUserID(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// Create ReadLatestLegalAgreementSigningByUserIDRequest struct from input JSON
	argBytes := []byte(args[0])
	var request ReadLatestLegalAgreementSigningByUserIDRequest
	if err := json.Unmarshal(argBytes, &request); err != nil {
		return shim.Error(fmt.Sprintf("Error unmarshaling ReadLatestLegalAgreementSigningByUserIDRequest: %s", err))
	}

	// Get iterator for all entries
	iterator, err := stub.GetStateByRange("", "")
	if err != nil {
		shim.Error(fmt.Sprintf("Error getting state iterator: %s", err))
	}
	defer iterator.Close()

	// Get the latest record
	var latestLegalAgreementSigning LegalAgreementSigning
	for iterator.HasNext() {
		// Get the next item
		item, err := iterator.Next()
		if err != nil {
			shim.Error(fmt.Sprintf("Error getting next item: %s", err))
		}

		// Unmarshal item
		var legalAgreementSigning LegalAgreementSigning
		err = json.Unmarshal(item.Value, &legalAgreementSigning)
		if err != nil && err.Error() != "Not a LegalAgreementSigning" {
			shim.Error(fmt.Sprintf("Error unmarshaling item: %s", err))
		}

		// Update latest record
		if latestLegalAgreementSigning.Timestamp < legalAgreementSigning.Timestamp &&
			legalAgreementSigning.UserID == request.UserID {
			latestLegalAgreementSigning = legalAgreementSigning
		}
	}

	// Return 404 if result's empty
	if len(latestLegalAgreementSigning.UserID) == 0 {
		return peer.Response{
			Status:  404,
			Message: fmt.Sprintf("Legal Agreement Signing for user %s does not exist", request.UserID),
		}
	}

	// Marshal latest record
	legalAgreementSigningAsBytes, _ := json.Marshal(latestLegalAgreementSigning)

	return shim.Success(legalAgreementSigningAsBytes)
}
