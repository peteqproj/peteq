package cmd

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/peteqproj/peteq/pkg/client"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/spf13/cobra"
)

var registerCmdFlags struct {
	email     string
	password  string
	autoLogin bool
	url       string
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register new user",
	RunE: func(cmd *cobra.Command, args []string) error {
		logr := logger.New(logger.Options{})
		u, err := url.Parse(registerCmdFlags.url)
		if err != nil {
			return err
		}
		c := client.NewAPIClient(&client.Configuration{
			HTTPClient: createHTTPClient(),
			Scheme:     u.Scheme,
			Servers: client.ServerConfigurations{
				{
					URL: u.Host,
				},
			},
		})
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
		token, err := login(c, 0, 10, logr)
		if err != nil {
			return err
		}
		fmt.Println(token)
		return storeClientConfiguration(registerCmdFlags.url, token)
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
	registerCmd.Flags().StringVar(&registerCmdFlags.email, "email", "", "Email")
	registerCmd.Flags().StringVar(&registerCmdFlags.password, "password", "", "Password")
	registerCmd.Flags().BoolVar(&registerCmdFlags.autoLogin, "login", true, "Also obtain authentication token after login")
	registerCmd.Flags().StringVar(&registerCmdFlags.url, "url", "http://localhost", "Peteq API url")

	registerCmd.MarkFlagRequired("email")
	registerCmd.MarkFlagRequired("password")
}

func login(c *client.APIClient, attempt int, maxTries int, logr logger.Logger) (string, error) {
	time.Sleep(time.Second * 5)
	logr.Info("Requesting API Token", "attempt", attempt, "max-retries", maxTries)
	data, response, err := c.UserCommandAPIApi.CUserLoginPost(context.Background()).Body(client.LoginRequestBody{
		Email:    registerCmdFlags.email,
		Password: registerCmdFlags.password,
	}).Execute()
	if response != nil && response.StatusCode == 400 {
		return login(c, attempt+1, maxTries, logr)
	}
	if err != nil {
		return "", err
	}

	if data.Status != nil && *data.Status == "rejected" {
		return "", fmt.Errorf("Failed to login: %s", *data.Reason)
	}
	return (*data.Data)["token"].(string), nil
}
