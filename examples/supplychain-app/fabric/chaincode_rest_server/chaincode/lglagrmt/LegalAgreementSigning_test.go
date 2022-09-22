package lglagrmt

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	. "github.com/chaincode/common"

	"github.com/franela/goblin"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	. "github.com/onsi/gomega"
)

// NewMockStub creates a MockStub. This currently requires using fabric builds from master branch
// as it requires the changes below, that are yet to be released: https://jira.hyperledger.org/browse/FAB-5644
func NewMockStub(name string, cc shim.Chaincode) *shim.MockStub {
	// Create new mock
	s := shim.NewMockStub(name, cc)
	return s
}

func TestLegalAgreementSigning(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	txID := "mockTxID"
	var mockStub *shim.MockStub
	chaincode := new(SmartContract)

	g.Describe("Init", func() {
		g.It("should initialize successfully", func() {
			mockStub = NewMockStub("mockstub", chaincode)

			mockStub.MockTransactionStart(txID)
			response := chaincode.Init(mockStub)
			chaincode.logger.SetLevel(shim.LogError)
			mockStub.MockTransactionEnd(txID)

			Expect(response.Status).To(BeEquivalentTo(200))
		})
	})

	g.Describe("Create Legal Agreement Signing", func() {
		g.BeforeEach(func() {
			mockStub = NewMockStub("mockstub", chaincode)
		})

		g.Describe("with valid data", func() {
			g.It("should return successfully", func() {
				// Read input fixture
				byteValue := readJSON(g, "../testdata/legal-agreement-signing-input-valid.json")
				var input LegalAgreementSigningRequest
				json.Unmarshal([]byte(byteValue), &input)

				// Run Create Legal Agreement Signing transaction
				args := [][]byte{[]byte("createLegalAgreementSigning"), byteValue}
				response := mockStub.MockInvoke("legalagreement", args)

				// Retrieve results
				var results map[string]interface{}
				json.Unmarshal(response.Payload, &results)

				Expect(response.Status).To(BeEquivalentTo(200))
				Expect(results["createdID"]).To(Equal(input.ID))
			})

			g.It("should write the legal agreement signing to the blockchain", func() {
				// Read input fixture
				byteValue := readJSON(g, "../testdata/legal-agreement-signing-input-valid.json")
				var input LegalAgreementSigningRequest
				json.Unmarshal([]byte(byteValue), &input)

				// Run Create Legal Agreement Signing transaction
				args := [][]byte{[]byte("createLegalAgreementSigning"), byteValue}
				response := mockStub.MockInvoke("legalagreement", args)

				Expect(response.Status).To(BeEquivalentTo(200))

				// Retrieve results from ledger
				bytes, _ := mockStub.GetState(input.ID)
				var results map[string]interface{}
				json.Unmarshal(bytes, &results)

				// Read output fixture
				byteValue = readJSON(g, "../testdata/legal-agreement-signing-output.json")
				var output map[string]interface{}
				json.Unmarshal([]byte(byteValue), &output)

				Expect(results).To(Equal(output))
			})

			g.It("duplicate creation", func() {
				// Read input fixture
				byteValue := readJSON(g, "../testdata/legal-agreement-signing-input-valid.json")
				var input LegalAgreementSigningRequest
				json.Unmarshal([]byte(byteValue), &input)

				// Run Create Legal Agreement Signing transaction
				args := [][]byte{[]byte("createLegalAgreementSigning"), byteValue}
				response1 := mockStub.MockInvoke("legalagreement", args)

				response2 := mockStub.MockInvoke("legalagreement", args)

				// Retrieve results
				var results map[string]interface{}
				json.Unmarshal(response1.Payload, &results)

				Expect(response2.Status).To(BeEquivalentTo(403))
				Expect(response2.Message).To(BeEquivalentTo("Legal Agreement Signing 001 already exists"))
			})
		})

		g.Describe("with invalid data", func() {
			g.It("should return an error if < 1 argument", func() {
				// Run Create Product transaction
				args := [][]byte{[]byte("createLegalAgreementSigning")}
				response := mockStub.MockInvoke("legalagreement", args)

				Expect(response.Status).To(BeEquivalentTo(500))
				Expect(response.Message).To(Equal("Incorrect number of arguments. Expecting 1"))
			})

			g.It("should return an error if > 1 argument", func() {
				// Run Create Product transaction
				args := [][]byte{[]byte("createLegalAgreementSigning"), []byte(""), []byte("")}
				response := mockStub.MockInvoke("legalagreement", args)

				Expect(response.Status).To(BeEquivalentTo(500))
				Expect(response.Message).To(Equal("Incorrect number of arguments. Expecting 1"))
			})
		})
	})

	g.Describe("Read Legal Agreement Signing", func() {
		g.BeforeEach(func() {
			mockStub = NewMockStub("mockstub", chaincode)
		})

		g.Describe("with valid data", func() {
			g.It("should return successfully", func() {
				// Read input fixture
				legalAgreementSigning := LegalAgreementSigning{
					ID:                        "0001",
					UserID:                    "001",
					LegalAgreementID:          "001",
					LegalAgreementContentHash: "5c23ff0895c77c61680680097fa64202e1f1864463e3f8e27a2bd3fc2c30592b",
					Accepted:                  true,
					Timestamp:                 1654027884,
				}
				bytes, _ := json.Marshal(legalAgreementSigning)
				key := legalAgreementSigning.ID
				mockStub.MockTransactionStart(txID)
				mockStub.PutState(key, bytes)
				mockStub.MockTransactionEnd(txID)

				// Run Read Legal Agreement Signing transaction
				args := [][]byte{[]byte("readLegalAgreementSigning"), []byte("0001")}
				response := mockStub.MockInvoke("legalagreement", args)

				// Retrieve results
				var result LegalAgreementSigning
				err := json.Unmarshal(response.Payload, &result)
				if err != nil {
					g.Fail(err)
				}

				Expect(response.Status).To(BeEquivalentTo(200))
				Expect(result).To(Equal(legalAgreementSigning))
			})
		})

		g.Describe("with invalid data", func() {
			g.It("should return an error if < 1 argument", func() {
				// Run Read Legal Agreement Signing transaction
				args := [][]byte{[]byte("readLegalAgreementSigning")}
				response := mockStub.MockInvoke("legalagreement", args)

				Expect(response.Status).To(BeEquivalentTo(500))
				Expect(response.Message).To(Equal("Incorrect number of arguments. Expecting 1"))
			})

			g.It("should return an error if > 1 argument", func() {
				// Run Read Legal Agreement Signing transaction
				args := [][]byte{[]byte("readLegalAgreementSigning"), []byte(""), []byte("")}
				response := mockStub.MockInvoke("legalagreement", args)

				Expect(response.Status).To(BeEquivalentTo(500))
				Expect(response.Message).To(Equal("Incorrect number of arguments. Expecting 1"))
			})

			g.It("should return 404 if the Legal Agreement Signing doesn't exist", func() {
				// Run Read Legal Agreement Signing transaction
				args := [][]byte{[]byte("readLegalAgreementSigning"), []byte("None")}
				response := mockStub.MockInvoke("legalagreement", args)

				Expect(response.Status).To(BeEquivalentTo(404))
				Expect(response.Message).To(Equal("Legal Agreement Signing None does not exist"))
			})
		})
	})

	g.Describe("Read Latest Legal Agreement Signing By User ID", func() {
		g.BeforeEach(func() {
			mockStub = NewMockStub("mockstub", chaincode)
		})

		g.Describe("with valid data", func() {
			g.It("should return successfully", func() {
				// Read input fixture
				legalAgreementSigning1 := LegalAgreementSigning{
					ID:                        "0001",
					UserID:                    "001",
					LegalAgreementID:          "001",
					LegalAgreementContentHash: "5c23ff0895c77c61680680097fa64202e1f1864463e3f8e27a2bd3fc2c30592b",
					Accepted:                  false,
					Timestamp:                 1654027884,
				}
				legalAgreementSigning2 := LegalAgreementSigning{
					ID:                        "0002",
					UserID:                    "001",
					LegalAgreementID:          "001",
					LegalAgreementContentHash: "5c23ff0895c77c61680680097fa64202e1f1864463e3f8e27a2bd3fc2c30592b",
					Accepted:                  true,
					Timestamp:                 1654027884,
				}
				legalAgreementSigningAsBytes1, _ := json.Marshal(legalAgreementSigning1)
				legalAgreementSigningAsBytes2, _ := json.Marshal(legalAgreementSigning2)
				key1 := legalAgreementSigning1.ID
				key2 := legalAgreementSigning2.ID
				mockStub.MockTransactionStart(txID)
				mockStub.PutState(key1, legalAgreementSigningAsBytes1)
				mockStub.PutState(key2, legalAgreementSigningAsBytes2)
				mockStub.MockTransactionEnd(txID)

				// Run Read Latest Legal Agreement Signing By User ID transaction
				args := [][]byte{[]byte("readLatestLegalAgreementSigningByUserID"), []byte("001")}
				response := mockStub.MockInvoke("legalagreement", args)

				// Retrieve results
				var result LegalAgreementSigning
				err := json.Unmarshal(response.Payload, &result)
				if err != nil {
					g.Fail(err)
				}

				Expect(response.Status).To(BeEquivalentTo(200))
				Expect(result).To(Equal(legalAgreementSigning2))
			})
		})

		g.Describe("with invalid data", func() {
			g.It("should return an error if < 1 argument", func() {
				// Run Read Latest Legal Agreement Signing By User ID transaction
				args := [][]byte{[]byte("readLatestLegalAgreementSigningByUserID")}
				response := mockStub.MockInvoke("legalagreement", args)

				Expect(response.Status).To(BeEquivalentTo(500))
				Expect(response.Message).To(Equal("Incorrect number of arguments. Expecting 1"))
			})

			g.It("should return an error if > 1 argument", func() {
				// Run Read Latest Legal Agreement Signing By User ID transaction
				args := [][]byte{[]byte("readLatestLegalAgreementSigningByUserID"), []byte(""), []byte("")}
				response := mockStub.MockInvoke("legalagreement", args)

				Expect(response.Status).To(BeEquivalentTo(500))
				Expect(response.Message).To(Equal("Incorrect number of arguments. Expecting 1"))
			})

			g.It("should return 404 if the Latest Legal Agreement Signing By User ID doesn't exist", func() {
				// Run Read Latest Legal Agreement Signing By User ID transaction
				args := [][]byte{[]byte("readLatestLegalAgreementSigningByUserID"), []byte("None")}
				response := mockStub.MockInvoke("legalagreement", args)

				Expect(response.Status).To(BeEquivalentTo(404))
				Expect(response.Message).To(Equal("Legal Agreement Signing for user None does not exist"))
			})
		})
	})
}
