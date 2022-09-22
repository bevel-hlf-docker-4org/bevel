package common

import (
	"encoding/json"
	"errors"
)

// LegalAgreementSigning stores signed legal agreements
type LegalAgreementSigning struct {
	ID                        string `json:"ID"`
	UserID                    string `json:"userID"`
	LegalAgreementID          string `json:"legalAgreementID"`
	LegalAgreementContentHash string `json:"legalAgreementContentHash"`
	Accepted                  bool   `json:"accepted"`
	Timestamp                 int64  `json:"timestamp"`
}

// UnmarshalJSON will override unmarshal
func (legalAgreementSigning *LegalAgreementSigning) UnmarshalJSON(data []byte) error {
	var input map[string]interface{}
	err := json.Unmarshal(data, &input)
	if err != nil {
		return err
	}

	if input["userID"] == nil {
		return errors.New("Not a LegalAgreementSigning")
	}

	// Prevent circular reference
	type Alias LegalAgreementSigning
	var output Alias
	err = json.Unmarshal(data, &output)
	if err != nil {
		return err
	}

	c := LegalAgreementSigning(output)
	*legalAgreementSigning = c

	return nil
}
