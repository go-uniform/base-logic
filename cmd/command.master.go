package cmd

import (
	"github.com/spf13/cobra"
	"service/service"
)

func init() {
	var firstName string
	var lastName string
	var email string
	var mobile string

	masterCmd := command("master", func(cmd *cobra.Command, args []string) {
		service.InitializeDiary(test, level, rate)
		service.Command("master", natsUri, compileNatsOptions(), map[string]string{})
	}, "Create a master admin user record if one does not already exist")

	masterCmd.Flags().StringVarP(&firstName, "first-name", "", "", "The master account's contact person's first name")
	masterCmd.Flags().StringVarP(&lastName, "last-name", "", "", "The master account's contact person's last name")
	masterCmd.Flags().StringVarP(&email, "email", "", "", "The master account's contact person email")
	masterCmd.Flags().StringVarP(&mobile, "mobile", "", "", "The master account's contact person mobile")

	if err := masterCmd.MarkFlagRequired("first-name"); err != nil {
		panic(err)
	}
	if err := masterCmd.MarkFlagRequired("last-name"); err != nil {
		panic(err)
	}
	if err := masterCmd.MarkFlagRequired("email"); err != nil {
		panic(err)
	}
	if err := masterCmd.MarkFlagRequired("mobile"); err != nil {
		panic(err)
	}

	rootCmd.AddCommand()
}