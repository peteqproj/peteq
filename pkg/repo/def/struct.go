package def

type (
	RepoDef struct {
		Name          string      `yaml:"name"`
		RootAggregate Aggregate   `yaml:"rootAggregate"`
		Aggregates    []Aggregate `yaml:"aggregates"`
		Database      Database    `yaml:"database"`
		// Tenant is a representation of tenant column
		// which if set must be part of all the sql queries
		// otherwise errNoTenantInContent is returned
		Tenant string `yaml:"tenant"`
	}
	Aggregate struct {
		Resource string `yaml:"resource"`
	}

	Database struct {
		Postgres PostgresDB `yaml:"postgres"`
	}

	PostgresDB struct {
		Columns []Column `yaml:"columns"`
		// Indexes will be used to create ListBy... Method
		Indexes [][]string `yaml:"indexes"`
		// UniqueIndexes will be used to create GetBy... Method
		UniqueIndexes [][]string `yaml:"uniqueIndexes"`
		PrimeryKey    []string   `yaml:"primeryKey"`
	}

	Column struct {
		Name string `yaml:"name"`
		Type string `yaml:"type"`
		From struct {
			Type string `yaml:"type"`
			Path string `yaml:"path"`
		} `yaml:"from"`
	}
)
