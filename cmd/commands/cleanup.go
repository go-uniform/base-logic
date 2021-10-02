package commands

import (
	"github.com/spf13/cobra"
	"service/cmd/_base"
	"service/service"
	"time"
)

func init() {
	cmd := _base.Command("cleanup", func(cmd *cobra.Command, args []string) {
		service.Command("cleanup", time.Second, _base.NatsUri, _base.CompileNatsOptions(), map[string]string{}, nil)
	}, "Execute the system data cleanup routine")

	_base.RootCmd.AddCommand(cmd)
}