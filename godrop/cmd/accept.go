package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

var acceptCmd = &cobra.Command{
	Use:   "accept [COMMAND]",
	Short: "Accept Files Shared by Peers",
	Run:   runAccept,
}

func runAccept(command *cobra.Command, args []string) {
	l := logger{
		command: command,
	}

	l.log("Configuring Godrop MDNS ...")
	drop, err := configGodropMdns()

	if err != nil {
		errorAndExit(err)
	}

	connStrategy := drop.NewP2PConn("mdns")

	l.log("Attempting to establish connection...")
	p2pConn, err := connStrategy.Connect("")

	if err != nil {
		errorAndExit(err)
	}

	l.log("Reading from peer...")

	f, err := os.Create("cuba.jpg")

	if err != nil {
		errorAndExit(err)
	}

	//writer := bufio.NewWriter(os.Stdout)
	data := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := p2pConn.Read(buf)

		if err != nil {
			if err == io.EOF {
				data = append(data, buf[:n]...)
				f.Write(data)
				f.Close()

				return
			}
		}

		data = append(data, buf[:n]...)
	}

}

func init() {
	RootCmd.AddCommand(acceptCmd)

}
