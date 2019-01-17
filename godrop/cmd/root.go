package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const config = ".godrop"

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

	// pathToConf := path.Join(home,config)
	// if _, err := os.Stat(pathToConf); err != nil {
	// 	//file does exist delete it first

	// }

	// if err := viper.Read
}
