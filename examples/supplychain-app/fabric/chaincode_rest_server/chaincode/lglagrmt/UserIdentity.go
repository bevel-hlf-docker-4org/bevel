package lglagrmt

import (
	"encoding/json"
	"fmt"

	. "github.com/chaincode/common"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// createUserIdentity creates an user identity in the ledger
func (s *SmartContract) createUserIdentity(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// Create UserIdentityRequest struct from input JSON
	argBytes := []byte(args[0])
	var request UserIdentityRequest
	if err := json.Unmarshal(argBytes, &request); err != nil {
		return shim.Error(fmt.Sprintf("Error unmarshaling UserIdentityRequest: %s", err))
	}

	// Check if user identity state using id as key exists
	testUserIdentityAsBytes, err := stub.GetState(request.UserID)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return 403 if item exists
	if len(testUserIdentityAsBytes) != 0 {
		return peer.Response{
			Status:  403,
			Message: fmt.Sprintf("User Identity %s already exists", request.UserID),
		}
	}

	// Create a new UserIdentity
	newUserIdentity := UserIdentity{
		UserID:                    request.UserID,
		LegalAgreementSigningTxID: request.LegalAgreementSigningTxID,
		VerifiableCredential:      request.VerifiableCredential,
		Status:                    request.Status,
	}

	// Marshal user identity
	userIdentityAsBytes, _ := json.Marshal(newUserIdentity)
	err = stub.PutState(newUserIdentity.UserID, userIdentityAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	response := map[string]interface{}{
		"createdID": newUserIdentity.UserID,
		"txID":      stub.GetTxID(),
	}
	bytes, _ := json.Marshal(response)

	s.logger.Infof("Wrote User Identity: %s\n", newUserIdentity.UserID)
	return shim.Success(bytes)
}

// readUserIdentity returns the user identity with the given id
func (s *SmartContract) readUserIdentity(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// Create ReadUserIdentityRequest struct from input JSON
	argBytes := []byte(args[0])
	var request ReadUserIdentityRequest
	if err := json.Unmarshal(argBytes, &request); err != nil {
		return shim.Error(fmt.Sprintf("Error unmarshaling ReadUserIdentityRequest: %s", err))
	}

	// Get the user identity state from the ledger
	userIdentityAsBytes, err := stub.GetState(request.UserID)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return 404 if user identity does not exist
	if len(userIdentityAsBytes) == 0 {
		return peer.Response{
			Status:  404,
			Message: fmt.Sprintf("User Identity %s does not exist", request.UserID),
		}
	}

	// Unmarshal user identity
	var userIdentity UserIdentity
	err = json.Unmarshal(userIdentityAsBytes, &userIdentity)
	if err != nil {
		return peer.Response{
			Status:  400,
			Message: fmt.Sprintf("Failed to unmarshal User Identity: %s", err),
		}
	}

	return shim.Success(userIdentityAsBytes)
}
