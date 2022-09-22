package common

// LegalAgreementSigningRequest models the request to create a legal agreement signing
type LegalAgreementSigningRequest struct {
	ID                        string `json:"ID"`
	UserID                    string `json:"userID"`
	LegalAgreementID          string `json:"legalAgreementID"`
	LegalAgreementContentHash string `json:"legalAgreementContentHash"`
	Accepted                  bool   `json:"accepted"`
	Timestamp                 int64  `json:"timestamp"`
}

// ReadLegalAgreementSigningRequest models the request to read a legal agreement signing
type ReadLegalAgreementSigningRequest struct {
	ID string `json:"ID"`
}

// ReadLatestLegalAgreementSigningByUserIDRequest models the request to read latest legal agreement signing by user ID
type ReadLatestLegalAgreementSigningByUserIDRequest struct {
	UserID string `json:"userID"`
}
