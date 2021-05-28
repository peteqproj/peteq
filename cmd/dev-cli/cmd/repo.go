package cmd

import (
	"bytes"
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/gertd/go-pluralize"
	"github.com/hairyhenderson/gomplate"
	"github.com/peteqproj/peteq/pkg/logger"
	repo "github.com/peteqproj/peteq/pkg/repo/def"
	"github.com/peteqproj/peteq/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

//go:embed templates/repo
var tmpl string
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
		r := &repo.RepoDef{}
		err = yaml.Unmarshal(d, r)
		utils.DieOnError(err, "")

		p := pluralize.NewClient()
		r.Root.Database.Postgres.DBName = p.Plural(r.Root.Database.Name)

		for i, a := range r.Aggregates {
			r.Aggregates[i].Database.Postgres.DBName = p.Plural(a.Database.Name)
		}

		dir := path.Join(wd, "domain", r.Name)
		err = os.MkdirAll(dir, os.ModePerm)
		logr.Info("Creating repo")
		funcs := gomplate.Funcs(nil)
		funcs["BuildInitQueries"] = buildInitQueries
		funcs["BuildIndexesFunction"] = buildIndexesFunction(r)
		funcs["BuildIndexesArgumentList"] = buildIndexesArgumentList
		funcs["BuildColumnVar"] = buildColumnVar(r)
		funcs["EmbedRepoDef"] = embedRepoDef
		res, err := templateRepo(funcs, r)
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

func templateRepo(funcs template.FuncMap, data interface{}) ([]byte, error) {
	out := new(bytes.Buffer)
	t := template.New("")
	t = t.Funcs(funcs)
	t, err := t.Parse(tmpl)
	if err != nil {
		return nil, err
	}
	if err := t.Execute(out, data); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func buildIndexesFunction(r *repo.RepoDef) func([]string) string {
	return func(indexes []string) string {
		fn := strings.Builder{}
		for _, i := range indexes {
			fn.WriteString(strings.Title(i))
		}
		return fn.String()
	}
}

func buildIndexesArgumentList(indexes []string, database repo.Database) string {
	list := []string{
		"ctx context.Context",
	}
	for _, i := range indexes {
		var col *repo.Column
		for _, c := range database.Postgres.Columns {
			if c.Name == i {
				col = &c
			}
			if col != nil {
				break
			}
		}
		if col == nil {
			panic("Column not found")
		}
		list = append(list, fmt.Sprintf("%s %s", i, postgresTypeToGolangType(col.Type)))
	}
	return strings.Join(list, ", ")
}

func postgresTypeToGolangType(t string) string {
	switch t {
	case "string", "json":
		return "string"
	default:
		return "string"
	}
}

func buildInitQueries(r repo.RepoDef) string {
	queries := []string{}
	queries = append(queries, createTableString(r.Root.Database.Postgres.DBName, r.Root.Database.Postgres.Columns, r.Root.Database.Postgres.PrimeryKey))
	for _, idx := range r.Root.Database.Postgres.Indexes {
		queries = append(queries, createIndexString(idx, false, r.Root.Database.Postgres.DBName))
	}

	for _, idx := range r.Root.Database.Postgres.UniqueIndexes {
		queries = append(queries, createIndexString(idx, true, r.Root.Database.Postgres.DBName))
	}
	res := strings.Builder{}
	res.WriteString("var queries = []string{\n")

	for _, aggregate := range r.Aggregates {
		queries = append(queries, createTableString(aggregate.Database.Postgres.DBName, aggregate.Database.Postgres.Columns, aggregate.Database.Postgres.PrimeryKey))
		for _, idx := range aggregate.Database.Postgres.UniqueIndexes {
			queries = append(queries, createIndexString(idx, true, aggregate.Database.Postgres.DBName))
		}
		for _, idx := range aggregate.Database.Postgres.Indexes {
			queries = append(queries, createIndexString(idx, false, aggregate.Database.Postgres.DBName))
		}
	}

	for _, q := range queries {
		res.WriteString(fmt.Sprintf("\t\"%s\",\n", q))
	}

	res.WriteString("}")
	return res.String()
}

func createTableString(table string, columns []repo.Column, pks []string) string {
	q := strings.Builder{}
	q.WriteString("CREATE TABLE IF NOT EXISTS ")
	q.WriteString(table)
	q.WriteString("( ")
	col := []string{}
	for _, c := range columns {
		col = append(col, fmt.Sprintf("%s %s not null", c.Name, c.Type))
	}
	col = append(col, fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(pks, ",")))
	q.WriteString(strings.Join(col, ","))
	q.WriteString(");")
	return q.String()
}

func createIndexString(idx []string, unique bool, db string) string {
	q := strings.Builder{}
	if unique {
		q.WriteString("CREATE UNIQUE INDEX IF NOT EXISTS ")
	} else {
		q.WriteString("CREATE INDEX IF NOT EXISTS ")
	}
	q.WriteString(fmt.Sprintf("%s ON %s ", strings.Join(idx, "_"), db))
	q.WriteString("( ")
	index := []string{}
	for _, i := range idx {
		index = append(index, i)
	}
	q.WriteString(strings.Join(index, ","))
	q.WriteString(");")
	return q.String()
}

func buildColumnVar(r *repo.RepoDef) func(repo.Column) string {
	return func(c repo.Column) string {
		if c.FromResource != nil {
			if c.FromResource.Path == "." {
				return fmt.Sprintf("table_column_%s, err := json.Marshal(resource)\n if err != nil {\n return err\n}", c.Name)
			}
			return fmt.Sprintf("table_column_%s := resource.%s", c.Name, c.FromResource.Path)
		}

		if c.FromTenant != nil {
			return fmt.Sprintf("table_column_%s := user.%s", c.Name, c.FromTenant.Path)
		}

		panic("From is not defined")
	}
}

func embedRepoDef(r *repo.RepoDef) string {
	b, err := yaml.Marshal(r)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("`%s`", string(b))
}
