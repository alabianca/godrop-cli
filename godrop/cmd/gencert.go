package cmd

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/spf13/viper"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var genCertCmd = &cobra.Command{
	Use:   "gencert [INSTANCE]",
	Short: "Generate a TLS Certificate",
	Run:   runCert,
}

func runCert(command *cobra.Command, args []string) {
	home, err := homedir.Dir()

	if err != nil {
		log.Fatal(err)
	}

	godropDir := path.Join(home, ".godrop")

	csrTemplate := certTemplate(viper.GetString("UID") + "." + viper.GetString("Host"))
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		log.Fatal(err)
	}

	csr, err := x509.CreateCertificateRequest(rand.Reader, csrTemplate, privKey)

	if err != nil {
		log.Fatal(err)
	}

	//securely store the private key
	keyFile, err := os.Create(path.Join(godropDir, "priv.pem"))

	if err != nil {
		log.Fatal(err)
	}

	pem.Encode(keyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	})

	// send the CERTIFICATE REQUEST
	pemBlock := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csr,
	})

	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", "http://104.248.183.179:80/csr", bytes.NewBuffer(pemBlock))

	if err != nil {
		log.Fatal(err)
	}

	res, err := httpClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	certFile, err := os.Create(path.Join(godropDir, "server.crt"))

	if err != nil {
		log.Fatal(err)
	}

	io.Copy(certFile, res.Body)
}

func init() {
	RootCmd.AddCommand(genCertCmd)
}

func certTemplate(cn string) *x509.CertificateRequest {
	return &x509.CertificateRequest{
		Subject:            pkix.Name{Organization: []string{"godrop IO"}, CommonName: cn},
		SignatureAlgorithm: x509.SHA256WithRSA,
	}
}
