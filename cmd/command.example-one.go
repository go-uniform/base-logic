package cmd

import (
	"github.com/nats-io/go-nats"
	"github.com/spf13/cobra"
	"go-uniform/base-logic/service"
)

var exampleOneCmd = &cobra.Command{
	Use:   "command:example-one",
	Short: "Request the running " + service.AppName + " to execute the example-one command",
	Long:  "Request the running " + service.AppName + " to execute the example-one command",
	Run: func(cmd *cobra.Command, args []string) {
		service.Command("example-one", natsUri)
	},
}

func init() {
	exampleOneCmd.Flags().StringVarP(&natsUri, "nats", "n", nats.DefaultURL, "The nats cluster URI")
	rootCmd.AddCommand(exampleOneCmd)
}