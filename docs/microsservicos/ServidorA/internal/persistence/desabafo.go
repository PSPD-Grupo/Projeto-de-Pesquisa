package persistence

import "time"

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
	result, err := Db.Exec(`INSERT INTO desabafo (texto, created_at) VALUES (?, ?)`, d.Texto, d.Created_at)
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
