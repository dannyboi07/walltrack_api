package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

var PrivateKey *rsa.PrivateKey
var PublicKey *rsa.PublicKey

func InitKeys() error {
	var (
		privateKeyBytes []byte
		err             error
	)
	privateKeyBytes, err = ioutil.ReadFile("private.pem")
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("Failed to read private.pem file, err: %s", err.Error())
		}

		fmt.Println("Generating RSA keys...")
		PrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return fmt.Errorf("Failed to private key, err: %s", err)
		}

		privateKeyBytes = x509.MarshalPKCS1PrivateKey(PrivateKey)

		var privateKeyBlock pem.Block = pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		}

		var privKeyFile *os.File
		privKeyFile, err = os.Create("private.pem")
		if err != nil {
			return fmt.Errorf("Failed to create private.pem file, err: %s", err)
		}

		fmt.Println("Encoding to private.pem file...")
		privateKeyBytes = pem.EncodeToMemory(&privateKeyBlock)
		if privateKeyBytes == nil {
			return fmt.Errorf("Failed to encode private pem block")
		}
		if _, err = privKeyFile.Write(privateKeyBytes); err != nil {
			return fmt.Errorf("Failed to write private pem block to file, err: %s", err)
		}

		fmt.Println("Private key written to file")

	}

	if privKeyPem, _ := pem.Decode(privateKeyBytes); privKeyPem != nil {
		PrivateKey, err = x509.ParsePKCS1PrivateKey(privKeyPem.Bytes)
		if err != nil {
			return fmt.Errorf("Failed to parse private.pem bytes, err: %s", err)
		}
	} else {
		return fmt.Errorf("Failed to decode private.pem file, err: %s", err)
	}

	PublicKey = &PrivateKey.PublicKey
	fmt.Println("Loaded RSA keys")

	return nil
}
