package persistence

import (
	"time"
)

type desabafoDb struct {
	Texto      string `sql:"texto"`
	Id         int    `sql:"id"`
	Created_at int64  `sql:"created_at"`
}

func NewDesabafoDb(texto string) *desabafoDb {
	retorno := desabafoDb{
		Texto:      texto,
		Created_at: time.Now().UnixMilli(),
	}
	return &retorno
}

func (d *desabafoDb) Insert() error {
<<<<<<< HEAD
	result, err := Db.Exec(`INSERT INTO desabafo (texto, created_at) VALUES (?, ?);`, d.Texto, d.Created_at)
=======
	result, err := Db.Exec(
		`INSERT INTO desabafo (id, texto, created_at)
		VALUES ((SELECT COALESCE(MAX(id), 0) + 1 FROM desabafo), ?, ?)`,
		d.Texto,
		d.Created_at,
	)
>>>>>>> 4551b46 (Server stub)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	d.Id = int(id)
	return nil
}

func (db *DataBase) GetNdesabafos(n int) ([]desabafoDb, error) {
<<<<<<< HEAD
	rows, err := db.db.Query(" SELECT * FROM desabafo ORDER BY created_at DESC LIMIT ?;", n)
=======
	rows, err := db.db.Query("SELECT id, texto, created_at FROM desabafo ORDER BY created_at DESC LIMIT ?", n)
>>>>>>> 4551b46 (Server stub)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []desabafoDb

	for rows.Next() {
		var d desabafoDb
		if err := rows.Scan(&d.Id, &d.Texto, &d.Created_at); err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil

}
