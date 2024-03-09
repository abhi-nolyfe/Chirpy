package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)


type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Increment the file server hits counter
			cfg.fileserverHits++

			// Call the next handler in the chain
			next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	htmlResponse := `
<html>
<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>
</html>
`
	w.Write([]byte(fmt.Sprintf(htmlResponse, cfg.fileserverHits)))
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	r := chi.NewRouter()
	r.Route("/app", func(r chi.Router) {
    r.Use(apiCfg.middlewareMetricsInc)
    fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
    r.Handle("/*", fsHandler)
})
	r.Get("/api/healthz", handlerReadiness)
	r.Get("/admin/metrics", apiCfg.getMetrics)
	r.Get("/api/reset", apiCfg.resetMetrics)
	r.Post("/api/validate_chirp", validate)

	corsMux := middlewareCors(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}