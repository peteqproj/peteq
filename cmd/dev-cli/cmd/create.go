package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var createCmdFlags struct {
	verbose bool
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create resources",
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			createCmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
}
