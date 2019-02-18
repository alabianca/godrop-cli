package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/alabianca/godrop"

	"github.com/spf13/viper"
)

func configGodropMdns(tlsDesired bool) (*godrop.Godrop, error) {

	options := make([]godrop.Option, 0)
	options = append(options, mdnsBasicConfig)

	if tlsDesired {
		options = append(options, mdnsTLSconfig)
	}

	drop, err := godrop.NewGodrop(options...)

	if err != nil {
		return nil, err
	}

	return drop, nil
}

func mdnsBasicConfig(drop *godrop.Godrop) {

	drop.Port = viper.GetInt("LocalPort")
	drop.Host = viper.GetString("Host")
	drop.UID = viper.GetString("UID")

}

func mdnsTLSconfig(drop *godrop.Godrop) {
	//load TLS root cert
	rootCrt, err := loadCertificate("root")

	if err != nil {
		fmt.Println(err)
		return
	}

	godropCrt, err := loadCertificate("server")

	if err != nil {
		fmt.Println(err)
		return
	}

	privKey, err := loadPrivateKey()

	if err != nil {
		fmt.Println(err)
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
