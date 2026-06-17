package connection

import (
	"ServidorA/internal/persistence"
	"encoding/json"
	"log"
	"net/http"
)

type FeedServer struct{}

// POST /desabafo
// Body JSON: {"texto": "..."}
func (fs *FeedServer) PostDesabafo(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Texto string `json:"texto"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "json inválido", http.StatusBadRequest)
		return
	}

	novo := persistence.NewDesabafoDb(body.Texto)
	if err := novo.Insert(); err != nil {
		log.Fatalf("could not insert: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id":         novo.Id,
		"texto":      novo.Texto,
		"created_at": novo.Created_at,
	})
}

func (fs *FeedServer) GetFeed(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Quant int `json:"quant"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Quant <= 0 {
		http.Error(w, "json inválido ou 'quant' ausente", http.StatusBadRequest)
		return
	}

	desabafos, err := persistence.Db.GetNdesabafos(body.Quant)
	if err != nil {
		log.Printf("could not get desabafos: %v", err)
		http.Error(w, "erro ao buscar desabafos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	flusher, canFlush := w.(http.Flusher)

	for _, d := range desabafos {
		if err := enc.Encode(map[string]any{
			"id":         d.Id,
			"texto":      d.Texto,
			"created_at": d.Created_at,
		}); err != nil {
			log.Println("erro ao enviar mensagem:", err)
			return
		}
		if canFlush {
			flusher.Flush() // envia cada linha imediatamente, como stream.Send()
		}
	}
}
