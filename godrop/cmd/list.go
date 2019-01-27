package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
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
	timeFlag = *listCmd.Flags().Int("time", 5, "How long to browse for")
	RootCmd.AddCommand(listCmd)
}

func ls(command *cobra.Command, args []string) {
	drop, err := configGodropMdns()

	if err != nil {
		log.Fatal(err)
	}

	entries, err := drop.Discover(time.Duration(timeFlag))

	if err != nil {
		log.Fatal(err)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.AlignRight)
	fmt.Fprintln(w, "Instance\tService\tDomain\tText")
	for _, entry := range entries {
		fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s", entry.Instance, entry.Service, entry.Domain, entry.Text))
	}
	w.Flush()
}
