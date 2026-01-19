package main

import (
	"log"
	"github.com/google/uuid"
	"time"
	"github.com/Legendary-Coder-GT/chirpy/internal/database"
	"github.com/Legendary-Coder-GT/chirpy/internal/auth"
	_ "github.com/lib/pq"
	"encoding/json"
	"net/http"
)

type requestBody struct {
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}
type errorBody struct {
	Body string `json:"error"`
}
type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string	`json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) jsonHandler(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %s", err)
		w.WriteHeader(500)
		return
	}
	user_id, err := auth.ValidateJWT(token, cfg.jwt_secret)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	decoder := json.NewDecoder(req.Body)
	params := requestBody{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	if len(params.Body) > 140 {
		res := errorBody{"Chirp is too long"}
		data, err := json.Marshal(res)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(data)
	} else {
		args := database.CreateChirpParams{params.Body, user_id}
		chirp, err := cfg.db.CreateChirp(req.Context(), args)
		if err != nil {
			log.Printf("Error creating chirp: %s", err)
			w.WriteHeader(500)
			return
		}
		res := Chirp{
			chirp.ID,
			chirp.CreatedAt,
			chirp.UpdatedAt,
			chirp.Body,
			chirp.UserID,
		}
		data, err := json.Marshal(res)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		w.Write(data)
	}
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		log.Printf("Error getting chirps: %s", err)
		w.WriteHeader(500)
		return
	}
	res := []Chirp{}
	for _, chirp := range chirps {
		new_chirp := Chirp{
			chirp.ID,
			chirp.CreatedAt,
			chirp.UpdatedAt,
			chirp.Body,
			chirp.UserID,
		}
		res = append(res, new_chirp)
	}
	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}

func (cfg *apiConfig) handlerSingleChirp(w http.ResponseWriter, req *http.Request) {
	id, _ := uuid.Parse(req.PathValue("chirpID"))
	chirp, err := cfg.db.GetChrip(req.Context(), id)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	res := Chirp{
		chirp.ID,
		chirp.CreatedAt,
		chirp.UpdatedAt,
		chirp.Body,
		chirp.UserID,
	}
	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}

func (cfg *apiConfig) handlerDeleteChirps(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %s", err)
		w.WriteHeader(401)
		return
	}
	jwt_user_id, err := auth.ValidateJWT(token, cfg.jwt_secret)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	id, _ := uuid.Parse(req.PathValue("chirpID"))
	chirp, err := cfg.db.GetChrip(req.Context(), id)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	if jwt_user_id != chirp.UserID {
		w.WriteHeader(403)
		return
	}
	err = cfg.db.DeleteChirp(req.Context(), id)
	if err != nil {
		log.Printf("Error deleting chirp: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(204)
}