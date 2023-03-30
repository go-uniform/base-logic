package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"service/cmd/_base"
	"service/service"
	"time"
)

func init() {
	var firstName string
	var lastName string
	var email string
	var mobile string
	var password string

	cmd := _base.Command("master", func(cmd *cobra.Command, args []string) {
		service.Command("master", time.Second, _base.NatsUri, _base.CompileNatsOptions(), map[string]interface{}{
			"firstName": firstName,
			"lastName":  lastName,
			"email":     email,
			"mobile":    mobile,
			"password":  password,
		}, func(bytes []byte) {
			fmt.Println(string(bytes))
			os.Exit(0)
		})
	}, "Create a master admin user record if one does not already exist")

	cmd.Flags().StringVarP(&firstName, "firstName", "", "", "The master account's contact person's first name")
	cmd.Flags().StringVarP(&lastName, "lastName", "", "", "The master account's contact person's last name")
	cmd.Flags().StringVarP(&email, "email", "", "", "The master account's contact person's email")
	cmd.Flags().StringVarP(&mobile, "mobile", "", "", "The master account's contact person's mobile")
	cmd.Flags().StringVarP(&password, "password", "", "", "The master account's contact person's password")

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
	if err := cmd.MarkFlagRequired("password"); err != nil {
		panic(err)
	}

	_base.RootCmd.AddCommand(cmd)
}
