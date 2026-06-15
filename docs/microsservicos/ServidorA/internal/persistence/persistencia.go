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
	getNewId() int 
	InsertNewDesabafo(desabafo desabafoDb)
	getNoDb() map[int64]desabafoDb
}

type DataBase struct {
	db *sql.DB
	no_db map[int64]desabafoDb
	cont int
}
var isNoDB bool = true
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
		no_db: make(map[int64]desabafoDb, 0),
		cont: 0,
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
}

func (d *DataBase) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.db.Exec(query, args...)
}

func (d *DataBase) getNewId() int {
	id := d.cont
	d.cont++;
	return id
}

func (d *DataBase) getNoDb() map[int64]desabafoDb {
	return d.no_db
}


