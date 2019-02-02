package cmd

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
)

func godropDir() (string, error) {
	home, err := homedir.Dir()

	if err != nil {
		return "", err
	}

	dropDir := path.Join(home, ".godrop")

	return dropDir, nil
}

func loadPublicKey() (*rsa.PublicKey, error) {
	dropDir, err := godropDir()

	if err != nil {
		return nil, err
	}

	publicPem, errPub := os.Open(path.Join(dropDir, "pub.pem"))
	pubBytes, err := ioutil.ReadAll(publicPem)

	if errPub != nil {
		return nil, errPub
	}

	pubBlock, pubRest := pem.Decode(pubBytes)

	if len(pubRest) > 0 {
		return nil, fmt.Errorf("Could not properly parse pub.pem")
	}

	key, err := x509.ParsePKCS1PublicKey(pubBlock.Bytes)

	if err != nil {
		return nil, err
	}

	return key, nil
}

func loadPrivateKey() (*rsa.PrivateKey, error) {
	dropDir, err := godropDir()

	if err != nil {
		return nil, err
	}

	privatePem, errPrv := os.Open(path.Join(dropDir, "priv.pem"))
	prvBytes, err := ioutil.ReadAll(privatePem)

	if errPrv != nil {
		return nil, errPrv
	}

	prvBlock, prvRest := pem.Decode(prvBytes)

	if len(prvRest) > 0 {
		return nil, fmt.Errorf("Could not properly parse pub.pem")
	}

	key, err := x509.ParsePKCS1PrivateKey(prvBlock.Bytes)

	if err != nil {
		return nil, err
	}

	return key, nil
}
