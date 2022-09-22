package common

import (
	"encoding/json"
	"errors"
)

// LegalAgreement stores legal agreements
type LegalAgreement struct {
	ID          string `json:"ID"`
	Content     string `json:"content"`
	ContentHash string `json:"hash"`
	Timestamp   int64  `json:"timestamp"`
	Version     int64  `json:"version"`
}

// UnmarshalJSON will override Unmarshal
func (legalAgreement *LegalAgreement) UnmarshalJSON(data []byte) error {
	var input map[string]interface{}
	err := json.Unmarshal(data, &input)
	if err != nil {
		return err
	}

	if input["version"] == nil {
		return errors.New("Not a LegalAgreement")
	}

	// Prevent circular reference
	type Alias LegalAgreement
	var output Alias
	err = json.Unmarshal(data, &output)
	if err != nil {
		return err
	}

	c := LegalAgreement(output)
	*legalAgreement = c

	return nil
}
