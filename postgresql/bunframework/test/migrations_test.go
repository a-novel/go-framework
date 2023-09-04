//go:build !race

package test

import (
	"context"
	"github.com/a-novel/go-framework/postgresql/bunframework"
	"github.com/a-novel/go-framework/test"
	"github.com/stretchr/testify/require"
	"io/fs"
	"testing"
)

func TestMigrations(t *testing.T) {
	db, sqlDB := test.GetPostgres(t, []fs.FS{})
	defer db.Close()
	defer sqlDB.Close()

	migrator := new(bunframework.MigrateConfig)
	migrator.RegisterSQLMigrations(migrations)

	require.NoError(t, migrator.Execute(context.TODO(), db))

	report := migrator.Report()
	require.NotNil(t, report)
	require.Len(t, report.Migrations.Applied(), 1)

	require.NoError(t, migrator.Execute(context.TODO(), db))

	report = migrator.Report()
	require.NotNil(t, report)
	require.Len(t, report.Migrations.Applied(), 0)

	require.NoError(t, migrator.Rollback(context.TODO(), db))
}
