package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/alabianca/godrop"

	"github.com/spf13/viper"
)

// Configure Godrop over MDNS
func getGodropConfigs() []godrop.Option {
	return []godrop.Option{
		mdnsBasicConfig,
		mdnsTLSconfig, // @todo: ignore this config if no tls is desired (ie; local testing)
	}
}

func configGodropMdns() (*godrop.Godrop, error) {
	drop, err := godrop.NewGodrop(getGodropConfigs()...)

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
