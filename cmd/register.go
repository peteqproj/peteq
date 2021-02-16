package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var registerCmdFlags struct {
	email     string
	password  string
	autoLogin bool
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register new user",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(registerCmdFlags)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringVar(&registerCmdFlags.email, "email", "", "Email")
	registerCmd.Flags().StringVar(&registerCmdFlags.password, "password", "", "Password")
	registerCmd.Flags().BoolVar(&registerCmdFlags.autoLogin, "login", true, "Also obtain authentication token after login")

	registerCmd.MarkFlagRequired("email")
	registerCmd.MarkFlagRequired("password")
}
