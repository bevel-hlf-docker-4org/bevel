package common

// UserIdentity stores user identities
type UserIdentity struct {
	UserID                    string `json:"userID"`
	LegalAgreementSigningTxID string `json:"legalAgreementSigningTxID"`
	VerifiableCredential      string `json:"verifiableCredential"`
	Status                    string `json:"status"`
}
