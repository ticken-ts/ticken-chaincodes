package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
	"os"
	"strconv"
	"ticken-event-contract/tickenevent"
)

const chaincodeName = "cc-event"
const chaincodeVersion = "1.0.0"

func main() {
	// add metadata and init transaction context
	tickenEventContract := new(tickenevent.Contract)
	tickenEventContract.Name = chaincodeName
	tickenEventContract.Info.Version = chaincodeVersion
	tickenEventContract.TransactionContextHandler = tickenevent.NewTransactionContext()

	cc, err := contractapi.NewChaincode(tickenEventContract)
	if err != nil {
		log.Panicf("error creating %s chaincode: %s", chaincodeName, err)
	}

	server := &shim.ChaincodeServer{
		CCID:     getEnvOrPanic("CHAINCODE_ID"),
		Address:  getEnvOrPanic("CHAINCODE_SERVER_ADDRESS"),
		CC:       cc,
		TLSProps: getTLSProperties(),
	}

	if err := server.Start(); err != nil {
		log.Panicf("error starting %s chaincode service: %s", chaincodeName, err)
	}
}

func getTLSProperties() shim.TLSProperties {
	// Check if chaincode is TLS enabled
	tlsDisabledStr := getEnvOrDefault("CHAINCODE_TLS_DISABLED", "true")
	key := getEnvOrDefault("CHAINCODE_TLS_KEY", "")
	cert := getEnvOrDefault("CHAINCODE_TLS_CERT", "")
	clientCACert := getEnvOrDefault("CHAINCODE_CLIENT_CA_CERT", "")

	// convert tlsDisabledStr to boolean
	tlsDisabled, _ := strconv.ParseBool(tlsDisabledStr)
	var keyBytes, certBytes, clientCACertBytes []byte
	var err error

	if !tlsDisabled {
		keyBytes, err = os.ReadFile(key)
		if err != nil {
			log.Panicf("error while reading the crypto file: %s", err)
		}
		certBytes, err = os.ReadFile(cert)
		if err != nil {
			log.Panicf("error while reading the crypto file: %s", err)
		}
	}
	// Did not request for the peer cert verification
	if clientCACert != "" {
		clientCACertBytes, err = os.ReadFile(clientCACert)
		if err != nil {
			log.Panicf("error while reading the crypto file: %s", err)
		}
	}

	return shim.TLSProperties{
		Disabled:      tlsDisabled,
		Key:           keyBytes,
		Cert:          certBytes,
		ClientCACerts: clientCACertBytes,
	}
}

func getEnvOrDefault(env, defaultVal string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		value = defaultVal
	}
	return value
}

func getEnvOrPanic(env string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		log.Panicf("required env value %s not foudn", env)
	}
	return value
}
