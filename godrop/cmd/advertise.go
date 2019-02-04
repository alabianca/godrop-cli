package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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

	//go mainLoop(server)

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
