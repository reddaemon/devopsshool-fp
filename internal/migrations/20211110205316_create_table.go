package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateTable, downCreateTable)
}

func upCreateTable(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS public.currency  (
			id serial4 NOT NULL,
			date_of_request date NULL,
			valute_id varchar(7) NULL,
			numcode int4 NULL,
			charcode varchar(3) NULL,
			nominal int4 NULL,
			value float8 NULL,
			"name" varchar(1024) NULL
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

func downCreateTable(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`
	DROP TABLE public.currency;
	`)
	if err != nil {
		return err
	}
	return nil
}
