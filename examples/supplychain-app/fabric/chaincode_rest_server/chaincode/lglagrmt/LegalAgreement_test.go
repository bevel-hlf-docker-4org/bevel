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

func readJSON(g *goblin.G, path string) []byte {
	jsonFile, err := os.Open(path)
	if err != nil {
		g.Fail(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue
}

func TestLegalAgreement(t *testing.T) {
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

	g.Describe("Create Legal Agreement", func() {
		g.BeforeEach(func() {
			mockStub = NewMockStub("mockstub", chaincode)
		})

		g.Describe("with valid data", func() {
			g.It("should return successfully", func() {
				// Read input fixture
				byteValue := readJSON(g, "../testdata/legal-agreement-input-valid.json")
				var input LegalAgreementRequest
				json.Unmarshal([]byte(byteValue), &input)

				// Run Create Legal Agreement transaction
				args := [][]byte{[]byte("createLegalAgreement"), byteValue}
				response := mockStub.MockInvoke("legalagreement", args)

				// Retrieve results
				var results map[string]interface{}
				json.Unmarshal(response.Payload, &results)

				Expect(response.Status).To(BeEquivalentTo(200))
				Expect(results["createdID"]).To(Equal(input.ID))
			})

			g.It("should write the Legal Agreement to the blockchain", func() {
				// Read input fixture
				byteValue := readJSON(g, "../testdata/legal-agreement-input-valid.json")
				var input LegalAgreementRequest
				json.Unmarshal([]byte(byteValue), &input)

				// Run Create Legal Agreement transaction
				args := [][]byte{[]byte("createLegalAgreement"), byteValue}
				response := mockStub.MockInvoke("legalagreement", args)

				Expect(response.Status).To(BeEquivalentTo(200))

				// Retrieve results from ledger
				bytes, _ := mockStub.GetState(input.ID)
				var results map[string]interface{}
				json.Unmarshal(bytes, &results)

				// Read output fixture
				byteValue = readJSON(g, "../testdata/legal-agreement-output.json")
				var output map[string]interface{}
				json.Unmarshal([]byte(byteValue), &output)

				Expect(results).To(Equal(output))
			})
		})

		g.Describe("with invalid data", func() {
			g.It("should return an error if < 1 argument", func() {
				// Run Create Legal Agreement transaction
				args := [][]byte{[]byte("createLegalAgreement")}
				response := mockStub.MockInvoke("legalagreement", args)

				Expect(response.Status).To(BeEquivalentTo(500))
				Expect(response.Message).To(Equal("Incorrect number of arguments. Expecting 1"))
			})

			g.It("should return an error if > 1 argument", func() {
				// Run Create Legal Agreement transaction
				args := [][]byte{[]byte("createLegalAgreement"), []byte(""), []byte("")}
				response := mockStub.MockInvoke("legalagreement", args)

				Expect(response.Status).To(BeEquivalentTo(500))
				Expect(response.Message).To(Equal("Incorrect number of arguments. Expecting 1"))
			})

			g.It("should return an error if the version is greater than the latest version", func() {
				// Read input fixture
				byteValue1 := readJSON(g, "../testdata/legal-agreement-input-valid.json")
				var input1 LegalAgreementRequest
				var input2 LegalAgreementRequest
				json.Unmarshal([]byte(byteValue1), &input1)
				json.Unmarshal([]byte(byteValue1), &input2)

				input2.ID = "002"

				// Run Create Legal Agreement transaction
				args1 := [][]byte{[]byte("createLegalAgreement"), byteValue1}
				response1 := mockStub.MockInvoke("legalagreement", args1)

				// Run Create Legal Agreement transaction with different ID
				byteValue2, _ := json.Marshal(input2)
				args2 := [][]byte{[]byte("createLegalAgreement"), string(byteValue2)}
				response2 := mockStub.MockInvoke("legalagreement", args2)

				// Retrieve results
				var results1 map[string]interface{}
				json.Unmarshal(response1.Payload, &results)

				var results2 map[string]interface{}
				json.Unmarshal(response2.Payload, &results)

				Expect(response1.Status).To(BeEquivalentTo(200))
				Expect(results1["createdID"]).To(Equal(input1.ID))
				Expect(response2.Status).To(BeEquivalentTo(500))
				Expect(response2.Message).To(Equal("The version 1 is not greater than the latest version 1"))
			})
		})
	})

	g.Describe("Read Legal Agreement", func() {
		g.BeforeEach(func() {
			mockStub = NewMockStub("mockstub", chaincode)
		})

		g.Describe("with valid data", func() {
			g.It("should return successfully", func() {
				// Read input fixture
				legalAgreement := LegalAgreement{
					ID:          "001",
					Content:     "some legal agreement content first version",
					ContentHash: "5c23ff0895c77c61680680097fa64202e1f1864463e3f8e27a2bd3fc2c30592b",
					Timestamp:   1654027884,
					Version:     1,
				}
				bytes, _ := json.Marshal(legalAgreement)
				key := legalAgreement.ID
				mockStub.MockTransactionStart(txID)
				mockStub.PutState(key, bytes)
				mockStub.MockTransactionEnd(txID)

				// Run Read Legal Agreement transaction
				args := [][]byte{[]byte("readLegalAgreement"), []byte("001")}
				response := mockStub.MockInvoke("legalagreement", args)

				// Retrieve results
				var result LegalAgreement
				err := json.Unmarshal(response.Payload, &result)
				if err != nil {
					g.Fail(err)
				}

				Expect(response.Status).To(BeEquivalentTo(200))
				Expect(result).To(Equal(legalAgreement))
			})
		})

		g.Describe("with invalid data", func() {
			g.It("should return an error if < 1 argument", func() {
				// Run Read Legal Agreement transaction
				args := [][]byte{[]byte("readLegalAgreement")}
				response := mockStub.MockInvoke("legalagreement", args)

				Expect(response.Status).To(BeEquivalentTo(500))
				Expect(response.Message).To(Equal("Incorrect number of arguments. Expecting 1"))
			})

			g.It("should return an error if > 1 argument", func() {
				// Run Read Legal Agreement transaction
				args := [][]byte{[]byte("readLegalAgreement"), []byte(""), []byte("")}
				response := mockStub.MockInvoke("legalagreement", args)

				Expect(response.Status).To(BeEquivalentTo(500))
				Expect(response.Message).To(Equal("Incorrect number of arguments. Expecting 1"))
			})

			g.It("should return 404 if the Legal Agreement doesn't exist", func() {
				// Run Read Legal Agreement transaction
				args := [][]byte{[]byte("readLegalAgreement"), []byte("None")}
				response := mockStub.MockInvoke("legalagreement", args)

				Expect(response.Status).To(BeEquivalentTo(404))
				Expect(response.Message).To(Equal("Legal Agreement None does not exist"))
			})
		})
	})

	g.Describe("Read Latest Version Legal Agreement", func() {
		g.BeforeEach(func() {
			mockStub = NewMockStub("mockstub", chaincode)
		})

		g.Describe("with valid data", func() {
			g.It("should return successfully", func() {
				// Read input fixture
				legalAgreement1 := LegalAgreement{
					ID:          "001",
					Content:     "some legal agreement content first version",
					ContentHash: "5c23ff0895c77c61680680097fa64202e1f1864463e3f8e27a2bd3fc2c30592b",
					Timestamp:   1654027884,
					Version:     1,
				}
				legalAgreement2 := LegalAgreement{
					ID:          "002",
					Content:     "some legal agreement content second version",
					ContentHash: "52af150fcae310d02e368906b05fe33a907d46e3121533675d31931850d4dba5",
					Timestamp:   1654028933,
					Version:     2,
				}
				legalAgreementAsBytes1, _ := json.Marshal(legalAgreement1)
				legalAgreementAsBytes2, _ := json.Marshal(legalAgreement2)
				key1 := legalAgreement1.ID
				key2 := legalAgreement2.ID
				mockStub.MockTransactionStart(txID)
				mockStub.PutState(key1, legalAgreementAsBytes1)
				mockStub.PutState(key2, legalAgreementAsBytes2)
				mockStub.MockTransactionEnd(txID)

				// Run Read Legal Agreement transaction
				args := [][]byte{[]byte("readLatestVersionLegalAgreement")}
				response := mockStub.MockInvoke("legalagreement", args)

				// Retrieve results
				var result LegalAgreement
				err := json.Unmarshal(response.Payload, &result)
				if err != nil {
					g.Fail(err)
				}

				Expect(response.Status).To(BeEquivalentTo(200))
				Expect(result).To(Equal(legalAgreement2))
				Expect(result.Content).To(Equal(legalAgreement2.Content))
				Expect(result.Version).To(BeEquivalentTo(legalAgreement2.Version))
			})
		})

		g.Describe("with invalid data", func() {
			g.It("should return an error if > 0 argument", func() {
				// Run Read Latest Version Legal Agreement transaction
				args := [][]byte{[]byte("readLatestVersionLegalAgreement"), []byte("")}
				response := mockStub.MockInvoke("legalagreement", args)

				Expect(response.Status).To(BeEquivalentTo(500))
				Expect(response.Message).To(Equal("Incorrect number of arguments. Expecting 0"))
			})
		})
	})
}
