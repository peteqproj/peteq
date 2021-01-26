package cmd

import (
	"os"

	"fmt"

	"github.com/spf13/cobra"
)

type (
	rootCmdOptions struct {
		verbose bool
	}
)

var rootOptions rootCmdOptions

var rootCmd = &cobra.Command{
	Use:     "peteq",
	Version: "0.1.0",
}

// Execute - execute the root command
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
