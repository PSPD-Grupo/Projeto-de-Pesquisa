package persistence

import (
	"testing"
)

func TestNewDesabafoDb(t *testing.T) {
	texto := "Teste de desabafo unitario"
	d := NewDesabafoDb(texto)

	if d.Texto != texto {
		t.Errorf("Esperava texto '%s', obteve '%s'", texto, d.Texto)
	}
	if d.Created_at == 0 {
		t.Error("Esperava Created_at diferente de zero")
	}
}

func TestInsertNoDB(t *testing.T) {
	// Salva estado original
	originalIsNoDB := isNoDB
	originalDb := Db
	defer func() {
		isNoDB = originalIsNoDB
		Db = originalDb
	}()

	// Configura para rodar em memoria
	isNoDB = true
	dbMock := &DataBase{
		no_db: make(map[int64]desabafoDb),
		cont:  0,
	}
	Db = dbMock

	d := NewDesabafoDb("Desabafo em memoria")
	err := d.Insert()
	if err != nil {
		t.Fatalf("Erro ao inserir: %v", err)
	}

	if d.Id != 0 {
		t.Errorf("Esperava ID 0 para o primeiro item, obteve %d", d.Id)
	}

	// Verifica se foi inserido no mapa
	saved, ok := dbMock.no_db[0]
	if !ok {
		t.Error("Desabafo nao foi salvo no mapa em memoria")
	}
	if saved.Texto != "Desabafo em memoria" {
		t.Errorf("Texto incorreto no mapa: %s", saved.Texto)
	}
}
