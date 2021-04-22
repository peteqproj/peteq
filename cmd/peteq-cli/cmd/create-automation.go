package cmd

import (
	"github.com/peteqproj/peteq/pkg/client"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/spf13/cobra"
)

var createAutomationFlags struct {
	Type string
}

var createAutomationCmd = &cobra.Command{
	Use:   "automation ...names",
	Short: "Create automation",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		logr := logger.New(logger.Options{})
		cnf, auth, err := createClientConfiguration()
		if err != nil {
			return err
		}
		c := client.NewAPIClient(cnf)

		for _, name := range args {
			resp, _, err := c.AutomationCommandAPIApi.CAutomationCreatePost(auth).Body(client.CreateAutomationRequestBody{
				Name: name,
				Type: createAutomationFlags.Type,
			}).Execute()
			if err != nil {
				return err
			}
			logr.Info("Automation created", "id", resp.Id)
		}
		return nil
	},
}

func init() {
	createCmd.AddCommand(createAutomationCmd)
	createAutomationCmd.Flags().StringVar(&createAutomationFlags.Type, "type", "", "Automation type")
}
