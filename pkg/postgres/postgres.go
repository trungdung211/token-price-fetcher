package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"

	"github.com/trungdung211/token-price-fetcher/internal/entities/model"
)

func NewPostgresDb(uri string, migrate bool) (db *bun.DB, err error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(uri)))
	db = bun.NewDB(sqldb, pgdialect.New())
	if viper.GetBool("debug") {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	// init tables
	initTable(db)

	// migrations
	if migrate {
		initMigrations(sqldb)
	}

	return
}

func initTable(db *bun.DB) {
	var ctx = context.Background()

	entities := []interface{}{
		new(model.UserConfig),
		new(model.Price),
	}

	for _, entity := range entities {
		_, err := db.NewCreateTable().Model(entity).IfNotExists().Exec(ctx)
		if err != nil {
			panic(fmt.Sprintf("CreateTable error: %v, entity: %T", err, entity))
		}
	}
}

func initMigrations(sqldb *sql.DB) {
	driver, err := postgres.WithInstance(sqldb, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./gen/migrations",
		"postgres", driver)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
}
