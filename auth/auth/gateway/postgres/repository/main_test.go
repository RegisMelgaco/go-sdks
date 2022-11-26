package repository

import (
	"os"
	"testing"

	"github.com/regismelgaco/go-sdks/auth/auth/gateway/postgres/migrate"
	"github.com/regismelgaco/go-sdks/postgres"
)

func TestMain(m *testing.M) {
	postgres.SetMigrationFunc(migrate.Migrate)
	teardown := postgres.SetupPgContainer()

	code := m.Run()

	teardown()

	os.Exit(code)
}
