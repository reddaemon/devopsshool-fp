package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upUsersTable, downUsersTable)
}

func upUsersTable(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`CREATE TABLE public.users (
			id serial4 NOT NULL,
			created_at timestamp NULL,
			name varchar(100) NULL,
			age int4 NULL,
			verify bool NULL,
			updated_at timestamp NULL DEFAULT now(),
			last_name varchar(100) NULL
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

func downUsersTable(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`DROP TABLE public.currency;`)
	if err != nil {
		return err
	}
	return nil
}
