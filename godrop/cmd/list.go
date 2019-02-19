package cmd

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "list all services in the local network",
	Run:   ls,
}

var timeFlag int

func init() {
	listCmd.Flags().IntVarP(&timeFlag, "time", "t", 5, "How long to browse for")
	RootCmd.AddCommand(listCmd)

}

func ls(command *cobra.Command, args []string) {
	drop, err := configGodropMdns(false)

	if err != nil {
		log.Fatal(err)
	}

	entries, err := drop.Discover(time.Duration(timeFlag))

	if err != nil {
		log.Fatal(err)
	}
	buf := new(bytes.Buffer)
	header := fmt.Sprintf("%-20s%-20s%-20s%-20s\n", "Instance", "Service", "Domain", "Text")
	buf.WriteString(header)
	for _, entry := range entries {
		line := fmt.Sprintf("%-20s%-20s%-20s%-20s\n", entry.Instance, entry.Service, entry.Domain, entry.Text)
		buf.WriteString(line)
	}
	fmt.Println(buf.String())
}
