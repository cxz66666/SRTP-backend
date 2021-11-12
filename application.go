package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	shell "github.com/ipfs/go-ipfs-api"
	// This package is needed so that all the preloaded plugins are loaded automatically
)

func main() {
	// hash, _ := uploadToIPFS("../../test1.txt")
	// downloadFromIPFS(hash)

	log.Println("============ application-golang starts ============")

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists("appUser") {
		err := populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"github.com",
		"fabric-samples",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)

	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}

	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	contract := network.GetContract("assetTransfer")

	log.Println("--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger")
	result, err := contract.SubmitTransaction("InitLedger")
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))

	result, err = contract.EvaluateTransaction("ReadAsset", "1")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println(string(result))
	log.Println("============ application-golang ends ============")

}

func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
	credPath := filepath.Join(
		"..",
		"github.com",
		"fabric-samples",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "User1@org1.example.com-cert.pem")
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}

	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}

	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))
	return wallet.Put("appUser", identity)
}

var sh *shell.Shell

func uploadToIPFS(filePath string) (string, error) {
	sh = shell.NewShell("localhost:5001")
	r, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	hash, err := sh.Add(r)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(hash)
	return hash, nil
}

func downloadFromIPFS(hash string) error {
	sh = shell.NewShell("localhost:5001")

	read, err := sh.Cat(hash)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(read)
	if err != nil {
		panic(err)
	}
	// print(body)
	err = ioutil.WriteFile(hash, body, 0666)
	if err != nil {
		panic(err)
	}
	print("Fetch success! File is store in ", hash)
	return nil
}

/*


export FABRICTOOLS=/home/ubuntu/fabric/scripts/fabric-samples/bin
export FABRIC_CFG_PATH=/home/ubuntu/fabric/scripts/fabric-samples/config
export PATH=$PATH:$FABRICTOOLS


export PATH=${PWD}/../github.com/fabric-samples/bin:$PATH
export FABRIC_CFG_PATH=${PWD}/../github.com/fabric-samples/config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/../github.com/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/../github.com/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

*/
