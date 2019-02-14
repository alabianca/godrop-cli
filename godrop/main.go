package main

import (
	"os"

	"github.com/alabianca/godrop-cli/godrop/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.WarnLevel) // only log warn level
	log.SetOutput(os.Stdout)

	cmd.RootCmd.Execute()

}
