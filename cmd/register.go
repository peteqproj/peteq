package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/peteqproj/peteq/pkg/client"
	"github.com/peteqproj/peteq/pkg/logger"
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
		logr := logger.New(logger.Options{})
		c := client.NewAPIClient(client.NewConfiguration())
		res, _, err := c.UserCommandAPIApi.CUserRegisterPost(context.Background()).Body(client.RegistrationRequestBody{
			Email:    registerCmdFlags.email,
			Password: registerCmdFlags.password,
		}).Execute()
		if err != nil {
			return err
		}
		if res.Status != nil && *res.Status == "rejected" {
			return fmt.Errorf("Failed to register: %s", *res.Reason)
		}
		logr.Info("Registred")
		if !registerCmdFlags.autoLogin {
			return nil
		}
		logr.Info("Requesting API token")
		time.Sleep(time.Second * 5)
		res, _, err = c.UserCommandAPIApi.CUserLoginPost(context.Background()).Body(client.LoginRequestBody{
			Email:    registerCmdFlags.email,
			Password: registerCmdFlags.password,
		}).Execute()
		if err != nil {
			return err
		}
		if res.Status != nil && *res.Status == "rejected" {
			return fmt.Errorf("Failed to login: %s", *res.Reason)
		}
		return storeClientConfiguration("http://localhost", (*res.Data)["token"].(string))
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
