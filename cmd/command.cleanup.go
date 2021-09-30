package cmd

import (
	"github.com/spf13/cobra"
	"service/service"
)

func init() {
	rootCmd.AddCommand(command("cleanup", func(cmd *cobra.Command, args []string) {
		service.InitializeDiary(test, level, rate)
		service.Command("cleanup", natsUri, compileNatsOptions(), map[string]string{})
	}, "Execute the system data cleanup routine"))
}