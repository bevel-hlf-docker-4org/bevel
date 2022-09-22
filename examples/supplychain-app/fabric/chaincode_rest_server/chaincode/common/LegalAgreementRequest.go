package common

// LegalAgreementRequest models the request to create a legal agreement
type LegalAgreementRequest struct {
	ID        string `json:"ID"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
	Version   int64  `json:"version"`
}

// ReadLegalAgreementRequest models the request to read a legal agreement
type ReadLegalAgreementRequest struct {
	ID string `json:"ID"`
}
