package cmd

import (
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create resource",
}

func init() {
	rootCmd.AddCommand(createCmd)
}
