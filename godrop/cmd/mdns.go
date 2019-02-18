package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/alabianca/godrop"

	"github.com/spf13/viper"
)

// Configure Godrop over MDNS

func configGodropMdns() (*godrop.Godrop, error) {
	drop, err := godrop.NewGodrop(mdnsConfig)

	if err != nil {
		return nil, err
	}

	return drop, nil
}

func mdnsConfig(drop *godrop.Godrop) {

	drop.Port = viper.GetInt("LocalPort")
	drop.Host = viper.GetString("Host")
	drop.UID = viper.GetString("UID")

	//load TLS
	rootCrt, err := loadCertificate("root")

	if err != nil {
		return
	}

	godropCrt, err := loadCertificate("server")

	if err != nil {
		return
	}

	privKey, err := loadPrivateKey()

	if err != nil {
		return
	}

	drop.RootCaCert = rootCrt
	drop.GodropCert = godropCrt
	drop.GodropPrivateKey = privKey
}

func strToInt(str string, size int) int64 {
	num, _ := strconv.ParseInt(str, 10, size)
	return num
}

func errorAndExit(err error) {
	fmt.Println(fmt.Errorf("%s", err))
	os.Exit(1)
}
