package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const description = `Uses quicktype (https://quicktype.io/) to generate Golang struct from JSON-Schema that represents an aggregate.
Make sure quicktype is installed on your machine: npm install -g quicktype`

var aggregateCmdFlags struct {
	schemaPath string
	name       string
	pkg        string
}

var aggregateCmd = &cobra.Command{
	Use:   "aggregate",
	Short: "Create aggregate",
	Long:  description,
	RunE: func(cmd *cobra.Command, args []string) error {
		logr := logger.New(logger.Options{})
		wd, err := os.Getwd()
		utils.DieOnError(err, "Failed to read current working dir")
		if err != nil {
			return err
		}
		dir := path.Join(wd, "domain", aggregateCmdFlags.pkg)
		err = os.MkdirAll(dir, os.ModePerm)
		utils.DieOnError(err, fmt.Sprintf("Failed to create directory: %s", dir))
		output := path.Join(dir, fmt.Sprintf("%s.go", aggregateCmdFlags.name))
		logr.Info("Creating aggregate", "output", output)
		err = run(output, aggregateCmdFlags.pkg, aggregateCmdFlags.schemaPath, logr)
		utils.DieOnError(err, "Failed to run quicktype command")
		return nil
	},
}

func init() {
	createCmd.AddCommand(aggregateCmd)
	aggregateCmd.Flags().StringVar(&aggregateCmdFlags.name, "name", "", "Aggregate name")
	aggregateCmd.Flags().StringVar(&aggregateCmdFlags.pkg, "package", "", "Package Name")
	aggregateCmd.Flags().StringVar(&aggregateCmdFlags.schemaPath, "schema", "", "Path to JSON-Schema")

	aggregateCmd.MarkFlagRequired("name")
	aggregateCmd.MarkFlagRequired("package")
	aggregateCmd.MarkFlagRequired("schema")

	aggregateCmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			aggregateCmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
}

func run(output string, pkg string, schema string, lgr logger.Logger) error {
	args := []string{
		"quicktype",
		"--lang", "go",
		"--src-lang", "schema",
		"-o", output,
		"--package", pkg,
		schema,
	}
	lgr.Info("Running quicktype", "cmd", strings.Join(args, " "))
	c := exec.Command("sh", "-c", strings.Join(args, " "))
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Start(); err != nil {
		return fmt.Errorf("Failed to run quicktype: %w", err)
	}
	return c.Wait()
}
