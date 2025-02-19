package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/tuanvumaihuynh/roboflow/cmd/roboflow_api/api"
	"github.com/tuanvumaihuynh/roboflow/internal/application"
	"github.com/tuanvumaihuynh/roboflow/pkg/cmdutils"
)

var rootCmd = &cobra.Command{
	Use:   "roboflow-standalone",
	Short: "Run Roboflow as a standalone single node server.",
	Run: func(_ *cobra.Command, _ []string) {
		run()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("error executing root command: %v", err)
	}
}

func run() {
	app, cleanup := application.New()
	defer func() {
		if err := cleanup(); err != nil {
			log.Fatalf("error cleaning up application: %v", err)
		}
	}()

	interruptChan := cmdutils.InterruptChan()

	go func() {
		if err := api.Start(app, interruptChan); err != nil {
			log.Fatalf("error starting api: %v", err)
		}
	}()

	<-interruptChan
}
