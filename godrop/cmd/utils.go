package cmd

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
)

const (
	GB = 1000000000 // 1 Gigabytes
	MB = 1000000    // 1 Megabyte
	KB = 1000       // 1 Byte
)

func godropDir() (string, error) {
	home, err := homedir.Dir()

	if err != nil {
		return "", err
	}

	dropDir := path.Join(home, ".godrop")

	return dropDir, nil
}

func loadCertificate(crt string) ([]byte, error) {
	godropDir, err := godropDir()

	if err != nil {
		return nil, err
	}

	root, err := ioutil.ReadFile(path.Join(godropDir, crt+".crt"))

	if err != nil {
		return nil, err
	}

	return root, nil

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

func loadPrivateKey() ([]byte, error) {
	dropDir, err := godropDir()

	if err != nil {
		return nil, err
	}

	privatePem, errPrv := os.Open(path.Join(dropDir, "priv.pem"))

	if errPrv != nil {
		return nil, errPrv
	}

	prvBytes, err := ioutil.ReadAll(privatePem)

	return prvBytes, nil

	// prvBlock, prvRest := pem.Decode(prvBytes)

	// if len(prvRest) > 0 {
	// 	return nil, fmt.Errorf("Could not properly parse priv.pem")
	// }

	// key, err := x509.ParsePKCS1PrivateKey(prvBlock.Bytes)

	// if err != nil {
	// 	return nil, err
	// }

	// return key, nil
}

func uint16ToBytes(num uint16) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, num); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func bytesToReadableFormat(numBytes int64) string {
	var res float32

	if numBytes >= GB {
		res = float32(numBytes) / float32(GB)
		return fmt.Sprintf("%.2f GB", res)
	}

	if numBytes >= MB {
		res = float32(numBytes) / float32(MB)
		return fmt.Sprintf("%.2f MB", res)
	}

	if numBytes >= KB {
		res = float32(numBytes) / float32(KB)
		return fmt.Sprintf("%.2f KB", res)
	}

	return fmt.Sprintf("%d Bytes", numBytes)
}
