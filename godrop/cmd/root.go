package cmd

import (
	"fmt"
	"net"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	config           = "godrop.yaml"
	defaultRelayIP   = "127.0.0.1"
	defaultRelayPort = "8080"
	defaultUID       = "def_godrop"
)

var RootCmd = &cobra.Command{
	Use:   "godrop",
	Short: "A P2P file share cli tool",
}

type logger struct {
	command *cobra.Command
}

func (l *logger) log(message string) {
	verbose := l.command.Flag("verbose").Value

	if verbose.String() == "true" {
		fmt.Println(message)
	}

}

func init() {
	cobra.OnInitialize(readConfig)
	RootCmd.PersistentFlags().Bool("verbose", false, "See logging")
}

func readConfig() {

	home, err := homedir.Dir()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	myIP, err := myIpv4()

	if err != nil {
		myIP = net.IPv4(byte(127), byte(0), byte(0), byte(1))
	}

	viper.AddConfigPath(home)
	viper.SetConfigName(config)
	viper.SetConfigType("yaml")

	// Check if the Config file exists. If not create it with defaults
	pathToConf := path.Join(home, ".godrop", config)
	viper.SetConfigFile(pathToConf)
	if _, err := os.Stat(pathToConf); os.IsNotExist(err) {
		//file does not exist create it first
		if err := os.Mkdir(path.Join(home, ".godrop"), 0700); err != nil {
			panic(err)
		}

		if _, err := os.Create(pathToConf); err != nil {
			panic(fmt.Errorf("Fatal error config file: %s\n", err))
		}

		viper.SetDefault("UID", defaultUID)
		viper.SetDefault("Host", "godrop.local")
		viper.SetDefault("LocalPort", 4000)
		viper.SetDefault("LocalIP", myIP.String())

		if e := viper.WriteConfig(); e != nil {
			panic(fmt.Errorf("Could not write config file %s\n", e))
		}

	}

	// Finally load in the file

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

}

func myIpv4() (net.IP, error) {
	ifaces, _ := net.Interfaces()

	for _, iface := range ifaces {

		addr, _ := iface.Addrs()

		for _, a := range addr {
			if strings.Contains(a.String(), ":") { //must be an ipv6
				continue
			}

			ip := addressStringToIP(a.String())

			if ip.IsLoopback() {
				continue
			}

			return ip, nil

		}
	}
	e := noIPv4FoundError{}
	return nil, e
}

func addressStringToIP(address string) net.IP {
	split := strings.Split(address, "/")
	ipSlice := strings.Split(split[0], ".")

	parts := make([]byte, 4)

	for i := range ipSlice {
		part, _ := strconv.ParseInt(ipSlice[i], 10, 16)
		parts[i] = byte(part)
	}

	ip := net.IPv4(parts[0], parts[1], parts[2], parts[3])

	return ip

}

type noIPv4FoundError struct{}

func (e noIPv4FoundError) Error() string {
	return "No IPv4 Interface found"
}
