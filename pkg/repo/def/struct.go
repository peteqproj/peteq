package def

type (
	RepoDef struct {
		Name string `yaml:"name"`

		// Tenant is a representation of tenant column
		// which if set must be part of all the sql queries
		// otherwise errNoTenantInContent is returned
		Tenant string `yaml:"tenant"`

		Root       Aggregate   `yaml:"root"`
		Aggregates []Aggregate `yaml:"aggregates"`
	}
	Aggregate struct {
		Resource string   `yaml:"resource"`
		Database Database `yaml:"database"`
	}

	Database struct {
		Name     string     `yaml:"name"`
		Postgres PostgresDB `yaml:"postgres"`
	}

	PostgresDB struct {
		// DBName calculated and will be overwrite
		DBName  string   `yaml:"dbname"`
		Columns []Column `yaml:"columns"`
		// Indexes will be used to create ListBy... Method
		Indexes [][]string `yaml:"indexes"`
		// UniqueIndexes will be used to create GetBy... Method
		UniqueIndexes [][]string `yaml:"uniqueIndexes"`
		PrimeryKey    []string   `yaml:"primeryKey"`
	}

	Column struct {
		Name         string `yaml:"name"`
		Type         string `yaml:"type"`
		FromResource *From  `yaml:"fromResource,omitempty"`
		FromTenant   *From  `yaml:"fromTenant,omitempty"`
	}

	From struct {
		As   string `yaml:"as"`
		Path string `yaml:"path"`
	}
)
