# OurGlass Chaincode for Legal Agreement

This is the **OurGlass** chaincode for the Legal Agreement smart contract. It is used to store the information of the Legal Agreement and Legal Agreement Signing.

> Keep in mind the following options for the next commands:

```bash
-c <args> # chaincode invocation arguments
-n <chaincode-name> # chaincode name
-C <channel-name> # channel name
```

## Transactions for the Legal Agreement

- [createLegalAgreement](#createlegalagreement)
- [readLegalAgreement](#readlegalagreement)
- [readLatestVersionLegalAgreement](#readlatestversionlegalagreement)

### createLegalAgreement

This transaction creates a new Legal Agreement. Run the following command to submit the transaction:

```bash
peer chaincode invoke -n <chaincode-name> -c '{"Args":["createLegalAgreement", "{\"ID\":\"001\",\"content\":\"some legal agreement content first version\",\"timestamp\":1653417608,\"version\":1}"]}' -C <channel-name>
```

### readLegalAgreement

This transaction reads the information of the Legal Agreement with the given ID. Run the following command to submit the transaction:

```bash
peer chaincode invoke -n <chaincode-name> -c '{"Args":["readLegalAgreement", "{\"ID\":\"001\"}"]}' -C <channel-name>
```

### readLatestVersionLegalAgreement

This transaction reads the information of the latest version of the Legal Agreement recorded in the ledger. Run the following command to submit the transaction:

```bash
peer chaincode invoke -n <chaincode-name> -c '{"Args":["readLatestVersionLegalAgreement"]}' -C <channel-name>
```

## Transactions for the Legal Agreement Signing

- [createLegalAgreementSigning](#createlegalagreementsigning)
- [readLegalAgreementSigning](#readlegalagreementsigning)
- [readLatestLegalAgreementSigningByUserID](#readlatestlegalagreementsigningbyuserid)

### createLegalAgreementSigning

This transaction creates a new Legal Agreement Signing. Run the following command to submit the transaction:

```bash
peer chaincode invoke -n <chaincode-name> -c '{"Args":["createLegalAgreementSigning", "{\"ID\":\"0001\",\"userID\":\"001\",\"legalAgreementID\":\"001\",\"legalAgreementContentHash\":\"5c23ff0895c77c61680680097fa64202e1f1864463e3f8e27a2bd3fc2c30592b\",\"accepted\":false,\"timestamp\":1653417620}"]}' -C <channel-name>
```

### readLegalAgreementSigning

This transaction reads the information of the Legal Agreement Signing with the given ID. Run the following command to submit the transaction:

```bash
peer chaincode invoke -n <chaincode-name> -c '{"Args":["readLegalAgreementSigning", "{\"ID\":\"0001\"}"]}' -C <channel-name>
```

### readLatestLegalAgreementSigningByUserID

This transaction reads the information of the latest Legal Agreement Signing recorded in the ledger by the given user ID. Run the following command to submit the transaction:

```bash
peer chaincode invoke -n <chaincode-name> -c '{"Args":["readLatestLegalAgreementSigningByUserID", "{\"userID\":\"001\"}"]}' -C <channel-name>
```

## Flow of the Smart Contract

The following command flow assumes you have a Hyperledger Fabric network and have chaincode installed, instantiated, and registered with the network.

First, create a Legal Agreement.

```bash
peer chaincode invoke -n legalagreement -c '{"Args":["createLegalAgreement", "{\"ID\":\"001\",\"content\":\"some legal agreement content first version\",\"timestamp\":1653417608,\"version\":1}"]}' -C myc
```

Then, you can read the information of the Legal Agreement.

```bash
peer chaincode invoke -n legalagreement -c '{"Args":["readLegalAgreement", "{\"ID\":\"001\"}"]}' -C myc
```

If you want to read the information of the latest version of the Legal Agreement, first you need another Legal Agreement.

```bash
peer chaincode invoke -n legalagreement -c '{"Args":["createLegalAgreement", "{\"ID\":\"002\",\"content\":\"some legal agreement content second version\",\"timestamp\":1653417708,\"version\":2}"]}' -C myc
```

Then, you can read the information of the latest version of the Legal Agreement.

```bash
peer chaincode invoke -n legalagreement -c '{"Args":["readLatestVersionLegalAgreement"]}' -C myc
```

After that, you can create a Legal Agreement Signing.

```bash
peer chaincode invoke -n legalagreement -c '{"Args":["createLegalAgreementSigning", "{\"ID\":\"0001\",\"userID\":\"001\",\"legalAgreementID\":\"002\",\"legalAgreementContentHash\":\"52af150fcae310d02e368906b05fe33a907d46e3121533675d31931850d4dba5\",\"accepted\":true,\"timestamp\":1653417620}"]}' -C myc
```

Then, you can read the information of the Legal Agreement Signing.

```bash
peer chaincode invoke -n legalagreement -c '{"Args":["readLegalAgreementSigning", "{\"ID\":\"0001\"}"]}' -C myc
```

And, you can read the information of the latest Legal Agreement Signing by the given user ID.

```bash
peer chaincode invoke -n legalagreement -c '{"Args":["readLatestLegalAgreementSigningByUserID", "{\"userID\":\"001\"}"]}' -C myc
```
