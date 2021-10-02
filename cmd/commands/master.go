package commands

import (
	"github.com/spf13/cobra"
	"service/cmd/_base"
	"service/service"
	"time"
)

func init() {
	var firstName string
	var lastName string
	var email string
	var mobile string

	cmd := _base.Command("master", func(cmd *cobra.Command, args []string) {
		service.Command("master", time.Second, _base.NatsUri, _base.CompileNatsOptions(), map[string]string{
			"firstName": firstName,
			"lastName": lastName,
			"email": email,
			"mobile": mobile,
		}, nil)
	}, "Create a master admin user record if one does not already exist")

	cmd.Flags().StringVarP(&firstName, "firstName", "", "", "The master account's contact person's first name")
	cmd.Flags().StringVarP(&lastName, "lastName", "", "", "The master account's contact person's last name")
	cmd.Flags().StringVarP(&email, "email", "", "", "The master account's contact person email")
	cmd.Flags().StringVarP(&mobile, "mobile", "", "", "The master account's contact person mobile")

	if err := cmd.MarkFlagRequired("firstName"); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired("lastName"); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired("email"); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired("mobile"); err != nil {
		panic(err)
	}

	_base.RootCmd.AddCommand(cmd)
}