package postgres

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/require"
)

var (
	db          *pgxpool.Pool
	migrateOnce sync.Once
	migrateFunc func(conn string) error
)

func SetMigrationFunc(fn func(conn string) error) {
	migrateOnce.Do(func() {
		migrateFunc = fn
	})
}

func SetupPgContainer() (teardown func()) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		panic(fmt.Sprintf("Could not connect to docker: %s", err))
	}

	resource, err := pool.Run("postgres", "13.3", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=test_db"})
	if err != nil {
		panic(err)
	}

	connConfig, err := pgxpool.ParseConfig(
		fmt.Sprintf(
			"postgres://postgres:postgres@localhost:%s/test_db?user=postgres&password=secret&sslmode=disable",
			resource.GetPort("5432/tcp"),
		),
	)
	if err != nil {
		panic(err)
	}

	connConfig.MaxConns = 1

	if err := pool.Retry(func() error {
		db, err = pgxpool.ConnectConfig(
			context.Background(),
			connConfig,
		)
		if err != nil {
			return err
		}
		return db.Ping(context.Background())
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	resource.Expire(20 * 60) // delete container after 20 mins

	return func() {
		err = pool.Purge(resource)
		if err != nil {
			panic(fmt.Sprintf("failed to purge resource: %s", err))
		}
	}
}

func GetDB(t *testing.T) *pgxpool.Pool {
	testName := strings.ReplaceAll(strings.ToLower(t.Name()), "/", "__")
	dbName := fmt.Sprintf("db_%v_%s", rand.Int(), testName)

	_, err := db.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s;", dbName))
	require.NoError(t, err, "failed to create new test db (name: %s): %w", dbName, err)

	newDBConnStr := fmt.Sprintf(
		"postgres://postgres:postgres@localhost:%d/%s?user=postgres&password=secret&sslmode=disable",
		db.Config().ConnConfig.Port,
		dbName,
	)

	newDB, err := pgxpool.Connect(context.Background(), newDBConnStr)
	require.NoError(t, err, "failed to create connection pool with new test db: %v | %s", err, newDBConnStr)

	err = newDB.Ping(context.Background())
	require.NoError(t, err, "failed to comunicate with new test db: %w", err)

	if migrateFunc != nil {
		err = migrateFunc(newDBConnStr)

		require.NoError(t, err, "failed to run migrations: %w", err)
	}

	return newDB
}
