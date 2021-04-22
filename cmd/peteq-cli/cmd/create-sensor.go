package cmd

import (
	"github.com/peteqproj/peteq/pkg/client"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/spf13/cobra"
)

var createSensorFlags struct {
}

var createSensorCmd = &cobra.Command{
	Use:   "sensor ...names",
	Short: "Create sensor",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		logr := logger.New(logger.Options{})
		cnf, auth, err := createClientConfiguration()
		if err != nil {
			return err
		}
		c := client.NewAPIClient(cnf)

		for _, name := range args {
			resp, _, err := c.SensorCommandAPIApi.CSensorCreatePost(auth).Body(client.SensorCreateRequestBody{
				Name: name,
			}).Execute()
			if err != nil {
				return err
			}
			logr.Info("Sensor created", "id", resp.Id)
		}
		return nil
	},
}

func init() {
	createCmd.AddCommand(createSensorCmd)
	// createSensorCmd.Flags().StringVar(&createSensorFlags.project, "project", "", "Assign the sensor to project")
}
