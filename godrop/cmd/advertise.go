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

func checkFile(fPath string) error {
	fInfo, err := os.Stat(fPath)

	if os.IsNotExist(err) {
		return err
	}

	if fInfo.IsDir() {
		return fmt.Errorf("A directory was provided")
	}

	return nil
}

func acceptConnections(server *godrop.Server) {

	for {
		sesh, err := server.Accept()

		if err != nil {
			log.Fatal("Could Not Accept Connection: ", err)
			continue
		}

		go transfer(sesh)

	}
}

func transfer(s *godrop.Session) {
	s.WriteHeader()
}
