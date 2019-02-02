package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var genkeyCmd = &cobra.Command{
	Use:   "genkey",
	Short: "generate a public/private key pair",
	Run:   genkey,
}

func genkey(command *cobra.Command, args []string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	home, err := homedir.Dir()

	if err != nil {
		fmt.Println("An error occured finding home directory")
		os.Exit(1)
	}

	godropDir := path.Join(home, ".godrop")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	publicKey := &privateKey.PublicKey

	privateBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	publicBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(publicKey),
	}

	publicOut, err := os.Create(path.Join(godropDir, "pub.pem"))

	if err != nil {
		fmt.Println("could not create pub.pem")
		os.Exit(1)
	}

	if err := pem.Encode(publicOut, publicBlock); err != nil {
		fmt.Println("Could not write out pub.pem")
	}

	privateOut, err := os.Create(path.Join(godropDir, "priv.pem"))

	if err != nil {
		fmt.Println("Could not create priv.pem")
		os.Exit(1)
	}

	if err := pem.Encode(privateOut, privateBlock); err != nil {
		fmt.Println("Could not wirte out priv.pem")
		os.Exit(1)
	}

	fmt.Println("Created public/private key pair")

}

func init() {
	RootCmd.AddCommand(genkeyCmd)
}
