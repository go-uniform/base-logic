package cmd

import (
	"github.com/spf13/cobra"
	"service/cmd/_base"
	service "service/service/_base"
)

func init() {
	var level string
	var rate int
	var limit int
	var test bool
	var env string
	// todo: add custom flag variables here

	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run " + service.AppName + " service",
		Long:  "Run " + service.AppName + " service",
		Run: func(cmd *cobra.Command, args []string) {
			service.InitializeDiary(test, level, rate)
			service.Execute(limit, test, _base.NatsUri, _base.CompileNatsOptions(), service.M{
				// todo: link custom flags to arg values here, example: "custom": custom,
			})
		},
	}

	// set the service's environment configurations via many command-line-interface (CLI) arguments
	runCmd.Flags().StringVarP(&level, "lvl", "l", "trace", "The logging level ['trace', 'debug', 'info', 'notice', 'warning', 'error', 'fatal'] that service is running in")
	runCmd.Flags().IntVarP(&rate, "rate", "r", 1000, "The sample rate of the trace logs used for performance auditing [set to -1 to log every trace]")
	runCmd.Flags().IntVarP(&limit, "limit", "x", 1000, "The messages per second that each topic worker will be limited to [set to 0 or less for maximum throughput]")
	runCmd.Flags().BoolVar(&test, "test", false, "A flag indicating if service should enter into test mode")
	runCmd.Flags().StringVarP(&env, "env", "e", "", "The environment ['prod', 'demo', 'staging', 'qa', 'dev', 'local'] that service is running in")
	// todo: add custom CLI flags here

	if err := runCmd.MarkFlagRequired("env"); err != nil {
		panic(err)
	}

	_base.RootCmd.AddCommand(runCmd)
}
