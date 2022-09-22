package lglagrmt

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// The SmartContract containing this chaincode
type SmartContract struct {
	logger *shim.ChaincodeLogger
}

// Init is called during chaincode instantiation to initialize any data.
func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	s.logger = shim.NewLogger("legalagreement")
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode.
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	function, args := stub.GetFunctionAndParameters()

	// Call the internal function based on the arguments supplied
	switch function {
	case "init":
		return s.Init(stub)
	case "createLegalAgreement":
		return s.createLegalAgreement(stub, args)
	case "readLegalAgreement":
		return s.readLegalAgreement(stub, args)
	case "readLatestVersionLegalAgreement":
		return s.readLatestVersionLegalAgreement(stub, args)
	case "createLegalAgreementSigning":
		return s.createLegalAgreementSigning(stub, args)
	case "readLegalAgreementSigning":
		return s.readLegalAgreementSigning(stub, args)
	case "readLatestLegalAgreementSigningByUserID":
		return s.readLatestLegalAgreementSigningByUserID(stub, args)
	case "createUserIdentity":
		return s.createUserIdentity(stub, args)
	case "readUserIdentity":
		return s.readUserIdentity(stub, args)
	default:
		fmt.Printf("Function for Invoke invalid or missing: %s, %s", function, args)
		return shim.Error(fmt.Sprintf("Function for Invoke invalid or missing: %s, %s", function, args))
	}
}
