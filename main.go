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
    numWallets := 10 // Number of wallets to generate
    for i := 0; i < numWallets; i++ {
        GenerateWallet(i)
    }
}

func GenerateWallet(index int) {

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

    projectDir, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }

    walletDir := projectDir + "/wallets/"
    if _, err := os.Stat(walletDir); os.IsNotExist(err) {
        os.Mkdir(walletDir, 0755)
    }

    file, _ := json.MarshalIndent(wallet, "", " ")
    _ = ioutil.WriteFile(fmt.Sprintf("%s%s_%d.json", walletDir, address, index), file, 0755)
    fmt.Println("Wallet generated: ", wallet)
}