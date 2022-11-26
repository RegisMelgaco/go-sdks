package migrate

import (
	"embed"
	"errors"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/regismelgaco/go-sdks/erring"
)

//go:embed *.sql
var migrationsFS embed.FS

func Migrate(connectionStr string) error {
	source, err := httpfs.New(http.FS(migrationsFS), ".")
	if err != nil {
		return erring.Wrap(err).
			Describe("failed to source migration files").
			Build()
	}

	m, err := migrate.NewWithSourceInstance("httpfs", source, connectionStr)
	if err != nil {
		return erring.Wrap(err).
			Describe("failed to instantiate migrate.Migrate").
			Build()
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return erring.Wrap(err).
			Describe("failed to migrate").
			Build()
	}

	return nil
}
