package main

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"

	"github.com/q231950/sputnik/cmd"
)

func main() {

	log.SetHandler(cli.New(os.Stderr))
	log.SetLevel(log.InfoLevel)

	log.Info("Starting 🛰  спутник")
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Errorf("%s", err)
		os.Exit(-1)
	}
}
