package persistence

import (
	"fmt"
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
		log.Printf("Executing Query for n = %d", n)
		rows, err := db.db.Query(fmt.Sprintf(`
			SELECT id, texto, created_at
			FROM desabafo
			ORDER BY created_at DESC
			LIMIT %d;
		`, n))
		if err != nil {
			log.Printf("erro na query: %v", err)
			return nil, err
		}
		defer rows.Close()
		
		count := 0
		for rows.Next() {
			count++
			var d desabafoDb
			err := rows.Scan(
				&d.Id,
				&d.Texto,
				&d.Created_at,
			)
			if err != nil {
				log.Printf("erro no scan em linha %d: %v", count, err)
				return nil, err
			}
			result = append(result, d)
		}
		log.Printf("Query loop finished. Scanned %d rows.", count)

		if err := rows.Err(); err != nil {
			log.Printf("rows.Err() erro detectado: %v", err)
			return nil, err
		}
	}

	return result, nil
}
