package postgres_test

import (
	"context"
	"testing"

	"github.com/regismelgaco/go-sdks/postgres"
	"github.com/stretchr/testify/assert"
)

func Test_SetupPgContainer_And_GetDB(t *testing.T) {
	t.Parallel()

	teardown := postgres.SetupPgContainer()
	t.Cleanup(teardown)

	db := postgres.GetDB(t)

	err := db.Ping(context.Background())
	assert.NoError(t, err)
}
