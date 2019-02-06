package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alabianca/godrop"

	"github.com/spf13/cobra"
)

var advertiseCmd = &cobra.Command{
	Use:   "advertise [COMMAND]",
	Short: "Accept Files Shared by Peers",
	Run:   runAdvertise,
}

func runAdvertise(command *cobra.Command, args []string) {
	drop, err := configGodropMdns()

	if err != nil {
		log.Fatal(err)
	}

	server, err := drop.NewMDNSService()

	if err != nil {
		log.Fatal(err)
	}

	go acceptConnections(server)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sig:
		log.Println("Shutdown...")
		server.Shutdown()
	}

}

func init() {
	RootCmd.AddCommand(advertiseCmd)

}

func acceptConnections(server *godrop.Server) {
	sesh, err := server.Accept()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Received a connection. Encryption: ", sesh.IsEncrypted())

	sesh.Write([]byte("Hello World"))
	sesh.Flush()
}
