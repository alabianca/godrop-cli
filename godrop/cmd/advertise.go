package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

const (
	BUF_SIZE = 1024
)

var advertiseCmd = &cobra.Command{
	Use:   "advertise [FILE_TO_SHARE]",
	Short: "Share a file with peers in the local network",
	Run:   runAdvertise,
}

func runAdvertise(command *cobra.Command, args []string) {

	if len(args) < 1 {
		command.Usage()
		os.Exit(1)
	}

	fPath := args[0]

	if err := checkFile(fPath); err != nil {
		log.Fatal(err)
	}

	drop, err := configGodropMdns()

	if err != nil {
		log.Fatal(err)
	}

	server, err := drop.NewMDNSService(fPath)

	go server.Start()

	defer server.Shutdown()

	if err != nil {
		log.Fatal(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sig:
		log.Println("Shutdown...")
	}

}

func init() {
	RootCmd.AddCommand(advertiseCmd)

}

func checkFile(fPath string) error {
	_, err := os.Stat(fPath)

	if os.IsNotExist(err) {
		return err
	}

	return nil
}
