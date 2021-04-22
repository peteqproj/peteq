package cmd

import (
	"github.com/peteqproj/peteq/pkg/client"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/spf13/cobra"
)

var createSensorAutomationBindingFlags struct {
	automation string
	sensor     string
}

var createSensorAutomationBindingCmd = &cobra.Command{
	Use:     "SensorAutomationBinding [name]",
	Aliases: []string{"sab"},
	Short:   "Create Sensor Automation Binding",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		logr := logger.New(logger.Options{})
		cnf, auth, err := createClientConfiguration()
		if err != nil {
			return err
		}
		c := client.NewAPIClient(cnf)
		resp, _, err := c.AutomationCommandAPIApi.CAutomationBindSensorPost(auth).Body(client.CreateSensorAutomationBindingRequestBody{
			Name:       args[0],
			Automation: createSensorAutomationBindingFlags.automation,
			Sensor:     createSensorAutomationBindingFlags.sensor,
		}).Execute()
		if err != nil {
			return err
		}
		logr.Info("SensorAutomationBinding created", "id", resp.Id)
		return nil
	},
}

func init() {
	createCmd.AddCommand(createSensorAutomationBindingCmd)
	createSensorAutomationBindingCmd.Flags().StringVar(&createSensorAutomationBindingFlags.automation, "automation", "", "Automation ID")
	createSensorAutomationBindingCmd.Flags().StringVar(&createSensorAutomationBindingFlags.sensor, "sensor", "", "Sensor ID")
}
