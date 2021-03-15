package cmd

import (
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts application componentes",
}

func init() {
	rootCmd.AddCommand(startCmd)
}
