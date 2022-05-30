package main

import (
	"database/sql"
	"fmt"

	_ "github.com/Dsmit05/metida/scripts/migrator/migrations"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

// Данный скрипт создаст таблицу test
func main() {
	connectURL := fmt.Sprintf("postgres://%v:%v@%v:%v/%v",
		"postgres",
		"postgres",
		"localhost",
		"5432",
		"metida",
	)
	conn, err := sql.Open("pgx", connectURL)
	if err != nil {
		fmt.Println("sql.Open() err:", err)
	}

	if err := goose.Up(conn, "."); err != nil {
		fmt.Println("goose.Up() err:", err)
	}
}
