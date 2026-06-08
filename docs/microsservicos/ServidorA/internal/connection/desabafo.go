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
		CreatedAt: novo_desabafo.Created_at,
	}
	return &retorno, nil
}

func (fs *FeedServer) GetFeed(fr *desabafo.FeedRequest, stream desabafo.Feed_GetFeedServer) error {
	desabafos, err := persistence.Db.GetNdesabafos(int(fr.Quant))
	if err != nil {
		log.Printf("could not get desabafos: %v", err)
		return err
	}
	for _, d := range desabafos {
		var d_retorno *desabafo.Desabafo
		d_retorno = &desabafo.Desabafo{
			Id:        int32(d.Id),
			Texto:     d.Texto,
			CreatedAt: d.Created_at,
		}
		if err = stream.Send(d_retorno); err != nil {
			log.Println("erro ao enviar mensagem:", err)
			return err
		}
	}
	return nil
}
