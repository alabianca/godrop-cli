package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/alabianca/godrop"

	"github.com/spf13/viper"
)

func configGodropMdns() (*godrop.Godrop, error) {
	drop, err := godrop.NewGodrop(mdnsConfig)

	if err != nil {
		return nil, err
	}

	return drop, nil
}

func mdnsConfig(drop *godrop.Godrop) {

	pubKey, puberr := loadPublicKey()
	prvKey, prverr := loadPrivateKey()

	if puberr != nil || prverr != nil {
		drop.PrvKey = nil
		drop.PubKey = nil
		log.Println("Warning: Could not load private/public key pair. Communication will not be encrypted")
	} else {
		drop.PrvKey = prvKey
		drop.PubKey = pubKey
	}

	drop.Port = viper.GetInt("LocalPort")
	drop.Host = viper.GetString("Host")
	drop.UID = viper.GetString("UID")
}

func strToInt(str string, size int) int64 {
	num, _ := strconv.ParseInt(str, 10, size)
	return num
}

func errorAndExit(err error) {
	fmt.Println(fmt.Errorf("%s", err))
	os.Exit(1)
}
