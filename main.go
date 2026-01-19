package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	_ "github.com/lib/pq"
	"database/sql"
	"github.com/Legendary-Coder-GT/chirpy/internal/database"
	"github.com/joho/godotenv"
	"os"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	jwt_secret string
	polka_key  string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	secretKey := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_KEY")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error opening database", err)
		return
	}
	const filepathRoot = "."
	const port = "8080"
	apiCfg := &apiConfig{fileserverHits: atomic.Int32{}, db: database.New(db), jwt_secret: secretKey, polka_key: polkaKey}
	sm := http.NewServeMux()
	sm.HandleFunc("GET /api/healthz", handlerReadiness)
	sm.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	sm.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	sm.HandleFunc("POST /api/chirps", apiCfg.jsonHandler)
	sm.HandleFunc("POST /api/users", apiCfg.userHandler)
	sm.HandleFunc("PUT /api/users", apiCfg.updateUserHandler)
	sm.HandleFunc("GET /api/chirps", apiCfg.handlerChirps)
	sm.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerSingleChirp)
	sm.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	sm.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	sm.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	sm.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirps)
	sm.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerUpgrade)
	handler := http.StripPrefix("/app" , http.FileServer(http.Dir(filepathRoot)))
	sm.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	local_server := &http.Server{
		Addr:	 ":" + port,
		Handler: sm,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(local_server.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	str := fmt.Sprintf(
		`<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
		</html>`,
		cfg.fileserverHits.Load())
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte(str))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	godotenv.Load()
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		w.WriteHeader(403)
		return
	}
	cfg.fileserverHits.Store(0)
	err := cfg.db.DeleteUsers(req.Context())
	if err != nil {
		log.Printf("Error deleting users: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("Server Hits Reset to 0 - All users cleared"))
}