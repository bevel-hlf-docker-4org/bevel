package common

// UserIdentityRequest models the request to create an user identity
type UserIdentityRequest struct {
	UserID                    string `json:"userID"`
	LegalAgreementSigningTxID string `json:"legalAgreementSigningTxID"`
	VerifiableCredential      string `json:"verifiableCredential"`
	Status                    string `json:"status"`
}

// ReadUserIdentityRequest models the request to read an user identity
type ReadUserIdentityRequest struct {
	UserID string `json:"userID"`
}
