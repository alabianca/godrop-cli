package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	config           = ".godrop.yaml"
	defaultRelayIP   = "127.0.0.1"
	defaultRelayPort = "8080"
	defaultUID       = "def_godrop"
)

var RootCmd = &cobra.Command{
	Use:   "godrop",
	Short: "A P2P file share cli tool",
}

func init() {
	cobra.OnInitialize(readConfig)
}

func readConfig() {
	fmt.Println("In read config")
	home, err := homedir.Dir()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.AddConfigPath(home)
	viper.SetConfigName(config)
	//viper.SetConfigType("yaml")
	viper.SetDefault("RelayIP", defaultRelayIP)
	viper.SetDefault("RelayPort", defaultRelayPort)
	viper.SetDefault("UID", defaultUID)

	// Check if the Config file exists. If not create it with defaults
	pathToConf := path.Join(home, config)
	viper.SetConfigFile(pathToConf)
	if _, err := os.Stat(pathToConf); os.IsNotExist(err) {
		//file does not exist create it first
		if _, err := os.Create(pathToConf); err != nil {
			panic(fmt.Errorf("Fatal error config file: %s\n", err))
		}
		e := viper.WriteConfig()

	}

	// Finally load in the file

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

}
