package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

type Wallet struct {
	PublicKey  string `json:"PublicKey"`
	PrivateKey string `json:"PrivateKey"`
}

func main() {
	GenerateWallet()
}

func GenerateWallet() {

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])

	wallet := Wallet{
		PublicKey:  address,
		PrivateKey: hexutil.Encode(privateKeyBytes)[2:],
	}

	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(dirname + "/wallets/"); os.IsNotExist(err) {
		os.Mkdir(dirname+"/wallets/", 0755)
	}

	file, _ := json.MarshalIndent(wallet, "", " ")
	_ = ioutil.WriteFile(dirname+"/wallets/"+address+".json", file, 0755)
	fmt.Println("Wallet generated: ", wallet)
}
