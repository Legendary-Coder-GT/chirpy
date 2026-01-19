package main

import (
	"log"
	"net/http"
	"encoding/json"
	_ "github.com/lib/pq"
	"github.com/Legendary-Coder-GT/chirpy/internal/database"
	"github.com/Legendary-Coder-GT/chirpy/internal/auth"
	"github.com/google/uuid"
	"time"
)

type emailBody struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID        	   uuid.UUID `json:"id"`
	CreatedAt 	   time.Time `json:"created_at"`
	UpdatedAt 	   time.Time `json:"updated_at"`
	Email     	   string    `json:"email"`
	Token		   string	 `json:"token"`
	RefreshToken   string	 `json:"refresh_token"`
	IsChirpyRed	   bool		 `json:"is_chirpy_red"`
}

type DataBody struct {
	UserID string `json:"user_id"`
}

type WebhookBody struct {
	Event	string	 `json:"event"`
	Data	DataBody `json:"data"`
}

func (cfg *apiConfig) userHandler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := emailBody{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	user, err := cfg.db.CreateUser(req.Context(), params.Email)
	if err != nil {
		log.Printf("Error creating user in database: %s", err)
		w.WriteHeader(500)
		return
	}
	pw, _ := auth.HashPassword(params.Password)
	pw_params := database.UpdatePasswordParams{pw, params.Email, user.ID}
	err = cfg.db.UpdatePassword(req.Context(), pw_params)
	if err != nil {
		log.Printf("Error updating password: %s", err)
		w.WriteHeader(500)
		return
	}
	local_user := User{user.ID, user.CreatedAt, user.UpdatedAt, user.Email, "", "", user.IsChirpyRed.Bool}
	data, err := json.Marshal(local_user)
	if err != nil {
		log.Printf("Error marshalling json: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(data)
}

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %s", err)
		w.WriteHeader(401)
		return
	}
	user_id, err := auth.ValidateJWT(token, cfg.jwt_secret)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	decoder := json.NewDecoder(req.Body)
	params := emailBody{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	new_hash, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		w.WriteHeader(500)
		return
	}
	new_params := database.UpdatePasswordParams{new_hash, params.Email, user_id}
	err = cfg.db.UpdatePassword(req.Context(), new_params)
	if err != nil {
		log.Printf("Error updating users database: %s", err)
		w.WriteHeader(500)
		return
	}
	user, err := cfg.db.GetUserByID(req.Context(), user_id)
	if err != nil {
		log.Printf("Error getting user: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	resp := User{user.ID, user.CreatedAt, user.UpdatedAt, user.Email, token, "", user.IsChirpyRed.Bool}
	data, _ := json.Marshal(resp)
	w.Write(data)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := emailBody{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	user, err := cfg.db.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		log.Printf("Error getting user: %s", err)
		w.WriteHeader(500)
		return
	}
	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Error checking password: %s", err)
		w.WriteHeader(500)
		return
	}
	if !match {
		w.WriteHeader(401)
		w.Write([]byte("Incorrect email or password"))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		expiresIn, _ := time.ParseDuration("1h")
		token, err := auth.MakeJWT(user.ID, cfg.jwt_secret, expiresIn)
		if err != nil {
			log.Printf("Error generating access token: %s", err)
			w.WriteHeader(500)
			return
		}
		r_token, _ := auth.MakeRefreshToken()
		rt_params := database.CreateRefreshTokenParams{r_token, user.ID}
		_, err = cfg.db.CreateRefreshToken(req.Context(), rt_params)
		if err != nil {
			log.Printf("Error generating refresh token: %s", err)
			w.WriteHeader(500)
			return
		}
		local_user := User{user.ID, user.CreatedAt, user.UpdatedAt, user.Email, token, r_token, user.IsChirpyRed.Bool}
		data, _ := json.Marshal(local_user)
		w.Write(data)
	}
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %s", err)
		w.WriteHeader(500)
		return
	}
	full_token, err := cfg.db.GetUserFromRefreshToken(req.Context(), token)
	if err != nil || full_token.RevokedAt.Valid {
		w.WriteHeader(401)
		return
	}
	w.WriteHeader(200)
	type TokenJSON struct {
		Token string `json:"token"`
	}
	expiresIn, _ := time.ParseDuration("1h")
	a_token, err := auth.MakeJWT(full_token.UserID, cfg.jwt_secret, expiresIn)
	if err != nil {
		log.Printf("Error generating access token: %s", err)
		w.WriteHeader(500)
		return
	}
	local_token := TokenJSON{a_token}
	data, _ := json.Marshal(local_token)
	w.Write(data)
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %s", err)
		w.WriteHeader(500)
		return
	}
	err = cfg.db.RevokeToken(req.Context(), token)
	if err != nil {
		log.Printf("Error revoking token: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(204)
}

func (cfg *apiConfig) handlerUpgrade(w http.ResponseWriter, req *http.Request) {
	key, err := auth.GetAPIKey(req.Header)
	if err != nil {
		log.Print(err)
		w.WriteHeader(401)
		return
	}
	if key != cfg.polka_key {
		w.WriteHeader(401)
		return
	}
	decoder := json.NewDecoder(req.Body)
	params := WebhookBody{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}
	user_id, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		log.Printf("Error parsing User ID: %s", err)
		w.WriteHeader(500)
		return
	}
	err = cfg.db.UpgradeUserByID(req.Context(), user_id)
	if err != nil {
		log.Printf("Error upgrading user to Red: %s", err)
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(204)
}