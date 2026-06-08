package persistence

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/glebarez/go-sqlite"
)

type DataBaseInterface interface {
	Exec(query string, args ...interface{}) (sql.Result, error)

	GetNdesabafos(n int) ([]desabafoDb, error)
}

type DataBase struct {
	db *sql.DB
}

var Db DataBaseInterface

func StartDataBase() {
	var err error

	exPath := os.Getenv("ROOT")

	conn, err := sql.Open("sqlite", "file:"+exPath+"/db/db.sqlite3?_foreign_keys=on&_busy_timeout=5000")
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.Exec("PRAGMA foreign_keys = ON")
	Db = &DataBase{
		db: conn,
	}

	query, err := os.ReadFile(exPath + "/db/create_db.sql")
	if err != nil {
		log.Println("erro lendo arquivo sql")
		panic(err)
	}
	if _, err = Db.Exec(string(query)); err != nil {
		log.Println("erro executando arquivo create sql")
		panic(err)
	}
	if err = corrigirIdsAntigos(conn); err != nil {
		log.Println("erro corrigindo ids antigos")
		panic(err)
	}
}

func (d *DataBase) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.db.Exec(query, args...)
}

func corrigirIdsAntigos(conn *sql.DB) error {
	tx, err := conn.Begin()
	if err != nil {
		return err
	}

	// Se ocorrer erro em qualquer ponto, desfaz a transação.
	defer tx.Rollback()

	var nextID int

	err = tx.QueryRow(
		"SELECT COALESCE(MAX(id), 0) + 1 FROM desabafo",
	).Scan(&nextID)

	if err != nil {
		return err
	}

	rows, err := tx.Query(
		"SELECT rowid FROM desabafo WHERE id IS NULL ORDER BY rowid",
	)
	if err != nil {
		return err
	}

	var rowIDs []int64

	for rows.Next() {
		var rowID int64

		if err := rows.Scan(&rowID); err != nil {
			rows.Close()
			return err
		}

		rowIDs = append(rowIDs, rowID)
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}

	rows.Close()

	for _, rowID := range rowIDs {
		_, err := tx.Exec(
			"UPDATE desabafo SET id = ? WHERE rowid = ?",
			nextID,
			rowID,
		)

		if err != nil {
			return err
		}

		nextID++
	}

	return tx.Commit()
}
