package cmd

import (
	"encoding/json"
	"io/ioutil"

	"github.com/peteqproj/peteq/pkg/client"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/spf13/cobra"
)

var triggerFlags struct {
	jsonFile string
	data     string
}

var triggerCmd = &cobra.Command{
	Use:   "trigger [id]",
	Short: "Trigger sensor",
	RunE: func(cmd *cobra.Command, args []string) error {
		logr := logger.New(logger.Options{})
		cnf, auth, err := createClientConfiguration()
		if err != nil {
			return err
		}
		c := client.NewAPIClient(cnf)
		data := make(map[string]interface{})
		if triggerFlags.data != "" {
			if err := json.Unmarshal([]byte(triggerFlags.data), &data); err != nil {
				return err
			}
		}

		if triggerFlags.jsonFile != "" {
			bytes, err := ioutil.ReadFile(triggerFlags.jsonFile)
			if err != nil {
				return err
			}
			if err := json.Unmarshal(bytes, &data); err != nil {
				return err
			}
		}
		if _, _, err := c.SensorCommandAPIApi.CSensorTriggerPost(auth).Body(client.SensorRunRequestBody{
			Data: &data,
			Id:   args[0],
		}).Execute(); err != nil {
			return err
		}
		logr.Info("OK")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(triggerCmd)
	triggerCmd.Flags().StringVar(&triggerFlags.jsonFile, "file", "", "JSON file with the input data")
	triggerCmd.Flags().StringVar(&triggerFlags.data, "data", "", "Input JSON raw data")
}
