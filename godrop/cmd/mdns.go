package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/alabianca/godrop"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var mdnsCmd = &cobra.Command{
	Use:   "mdns [args]",
	Short: "share a file via mdns",
	Run:   runShare,
}

func runShare(command *cobra.Command, args []string) {
	if len(args) <= 0 {
		command.Usage()
		os.Exit(1)
	}

	l := logger{
		command: command,
	}

	path := args[0]

	fInfo, err := os.Stat(path)

	if os.IsNotExist(err) {
		fmt.Println(fmt.Errorf("Error %s\n", err))
		os.Exit(1)
	}

	fmt.Printf("Sharing %d bytes\n", fInfo.Size())

	share(path, &l)

}

func init() {
	shareCmd.AddCommand(mdnsCmd)

}

func share(path string, l *logger) {
	l.log("Configuring godrop via mdns...")
	drop, err := configGodropMdns()

	l.log("Configured Godrop MDNS")
	l.log(drop.Port)
	l.log(drop.IP)
	l.log(drop.Host)
	l.log(drop.ServiceName)

	if err != nil {
		errorAndExit(err)
	}

	connStrategy := drop.NewP2PConn("mdns")
	l.log("Trying to find a peer...")

	p2pConn, err := connStrategy.Connect("")

	l.log("Found a peer and connected")

	if err != nil {
		errorAndExit(err)
	}

	l.log("Reading file...")
	content, err := ioutil.ReadFile(path)

	if err != nil {
		errorAndExit(err)
	}

	l.log("Writing file...")
	p2pConn.Write(content)

}

func configGodropMdns() (*godrop.Godrop, error) {
	drop, err := godrop.NewGodrop()

	if err != nil {
		return nil, err
	}

	return drop, nil
}

func mdnsConfig(drop *godrop.Godrop) {
	drop.Port = viper.GetString("LocalPort")
	drop.IP = viper.GetString("LocalIP")
	drop.ServiceName = viper.GetString("ServiceName")
	drop.ServiceWeight = uint16(strToInt(viper.GetString("ServiceWeight"), 16))
	drop.TTL = uint32(strToInt(viper.GetString("TTL"), 32))
	drop.Priority = uint16(strToInt(viper.GetString("TTL"), 16))
	drop.Host = viper.GetString("Host")
}

func strToInt(str string, size int) int64 {
	num, _ := strconv.ParseInt(str, 10, size)
	return num
}

func errorAndExit(err error) {
	fmt.Println(fmt.Errorf("%s"), err)
	os.Exit(1)
}
