package persistence

import (
	"log"
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

func (d *DataBase) InsertNewDesabafo(desabafo desabafoDb) {
	d.no_db[int64(desabafo.Id)] = desabafo
}

func (d *desabafoDb) Insert() error {
	if isNoDB {
		d.Id = Db.getNewId()
		Db.InsertNewDesabafo(*d)
	}

	if !isNoDB {
		result, err := Db.Exec(`INSERT INTO desabafo (texto, created_at) VALUES (?, ?);`, d.Texto, d.Created_at)
		if err != nil {
			return err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		d.Id = int(id)
	}
	
	return nil
}

func (db *DataBase) GetNdesabafos(n int) ([]desabafoDb, error) {
	var result []desabafoDb

	if isNoDB {
		for _, value := range Db.getNoDb() {
			result = append(result, value)
		}
	}
	if !isNoDB {
		log.Printf("n: %d", n)
		rows, err := db.db.Query(`
			SELECT id, texto, created_at
			FROM Desabafo
			ORDER BY created_at DESC
			LIMIT ?;
		`, n)
		if err != nil {
			log.Printf("erro na query: %v", err)
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
		var d desabafoDb

		err := rows.Scan(
			&d.Id,
			&d.Texto,
			&d.Created_at,
		)
		if err != nil {
			log.Printf("erro no scan: %v", err)
			return nil, err
		}

		result = append(result, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	}

	return result, nil
}
