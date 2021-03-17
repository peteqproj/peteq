package cmd

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"path"

	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type (
	repo struct {
		Name          string      `yaml:"name"`
		RootAggregate Aggregate   `yaml:"rootAggregate"`
		Aggregates    []Aggregate `yaml:"aggregates"`
	}
	Aggregate struct {
		Resource string `yaml:"resource"`
	}
)

var repoCmdFlags struct {
	repo string
}
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Create repo",
	Long:  `Generate repository with access to all repos of it.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logr := logger.New(logger.Options{})
		wd, err := os.Getwd()
		utils.DieOnError(err, "Failed to read current working dir")

		d, err := ioutil.ReadFile(repoCmdFlags.repo)
		utils.DieOnError(err, "Failed to read file")

		r := &repo{}
		err = yaml.Unmarshal(d, r)
		utils.DieOnError(err, "")

		dir := path.Join(wd, "domain", r.Name)
		err = os.MkdirAll(dir, os.ModePerm)
		logr.Info("Creating repo")
		res, err := templateRepo(r)
		utils.DieOnError(err, "Failed to template repositry")
		err = ioutil.WriteFile(path.Join(dir, "repo.go"), res, os.ModePerm)
		utils.DieOnError(err, "Failed to write repo to file")
		return nil
	},
}

func init() {
	createCmd.AddCommand(repoCmd)
	repoCmd.Flags().StringVar(&repoCmdFlags.repo, "repo", "", "Path to repo.yaml")
	repoCmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			repoCmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
}

func templateRepo(data interface{}) ([]byte, error) {
	out := new(bytes.Buffer)
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return nil, err
	}
	if err := t.Execute(out, data); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

var tmpl = `
package {{ .Name }}

import (
	"context"
	"errors"
	"encoding/json"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/imdario/mergo"
	"github.com/peteqproj/peteq/pkg/db"
	"github.com/peteqproj/peteq/pkg/logger"
)

const db_name = "repo_{{.Name}}"

var errNotFound = errors.New("Resource not found")

type (
	Repo struct {
		DB db.Database 
		Logger logger.Logger
	}

	ListOptions struct {}
	GetOptions struct {
		ID    string
		Query string
	}
)

func (r *Repo) List(ctx context.Context, options ListOptions) ([]*{{ .RootAggregate.Resource }}, error) {
	return nil, nil
}

func (r *Repo) Get(ctx context.Context, options GetOptions) (*{{ .RootAggregate.Resource }}, error) {
	return nil, nil
}

func (r *Repo) Create(ctx context.Context, resource *{{ .RootAggregate.Resource }}) (error) {
	return nil
}

func (r *Repo) Delete(ctx context.Context, id string) (error) {
	return nil
}

func (r *Repo) Update(ctx context.Context, resource *{{ .RootAggregate.Resource }}) (error) {
	return nil
}
`
