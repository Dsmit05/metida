package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(Up, Down)
}

const add_table_test = `
CREATE TABLE test
(
    id          serial PRIMARY KEY,
    name        text
);
`

func Up(tx *sql.Tx) error {
	_, err := tx.Exec(add_table_test)
	if err != nil {
		return err
	}
	return nil
}

const drop_table_test = `
DROP TABLE test;
`

func Down(tx *sql.Tx) error {
	_, err := tx.Exec(drop_table_test)
	if err != nil {
		return err
	}
	return nil
}
