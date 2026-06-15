package connection

import (
	"ServidorA/gen/desabafo"
	"context"
	"log"

	"ServidorA/internal/persistence"
)

type FeedServer struct {
	desabafo.UnimplementedFeedServer
}

func (fs *FeedServer) PostDesabafo(_ context.Context, in *desabafo.RascunhoDesabafo) (*desabafo.Desabafo, error) {
	novo_desabafo := persistence.NewDesabafoDb(in.Texto)
	err := novo_desabafo.Insert()
	if err != nil {
		log.Fatalf("could not insert: %v", err)
	}

	retorno := desabafo.Desabafo{
		Id:        int32(novo_desabafo.Id),
		Texto:     novo_desabafo.Texto,
		CreatedAt: novo_desabafo.Created_at.UnixMilli(),
	}
	return &retorno, nil
}

