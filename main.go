package main

import (
    "net/http"
    "log"
    "web_server/internal/database"
    "os"
    "flag"
)

import _ "github.com/joho/godotenv/autoload"

type ApiConfig struct {
    fileserverHits int
    db          *database.DB
    jwtSecret   string
    polkaApiKey string
}

func main() {
    dbg := flag.Bool("debug", false, "Enable debug")
    flag.Parse()

    if *dbg {
        log.Print("Starting in Debug Mode")
        deleteDataBase("database.json")
    }

    db, err := database.NewDB("database.json")
    if err != nil {
        log.Fatal(err)
    }

    jwtSecret, ok := os.LookupEnv("JWT_SECRET")
    polkaApiKey, ok := os.LookupEnv("POLKA_API_KEY")
    if !ok {
        log.Fatal("Secret Could Not Be Loaded")
    }

    cfg := &ApiConfig{
        fileserverHits: 0,
        db:             db,
        jwtSecret:      jwtSecret,
        polkaApiKey:    polkaApiKey,
    }

    log.Printf("%v", cfg.jwtSecret)
    port := "8080"
    mux := http.NewServeMux()
    mux.Handle("/app/*", http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
    // mux.HandleFunc("GET /api/healthz", ReadinessHandler)
    mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
    mux.HandleFunc("GET /api/reset", cfg.resetHandler)
    mux.HandleFunc("POST /api/chirps", cfg.POSTChirpsHandler)
    mux.HandleFunc("GET /api/chirps", cfg.GETChirpsHandler)
    mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.GETChirpsIDHandler)
    mux.HandleFunc("POST /api/users", cfg.CreateUserHandler)
    mux.HandleFunc("POST /api/login", cfg.loginHandler)
    mux.HandleFunc("PUT /api/users", cfg.handlerUpdateUserInformation)
    mux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)
    mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)
    mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.HandlerDeleteChirpId)
    mux.HandleFunc("POST /api/polka/webhooks", cfg.handlerUpgradeToChirpyRed)

    svr := http.Server {
        Addr:   ":" + port,
        Handler:    mux,
    }
    log.Fatal(svr.ListenAndServe())
}

func deleteDataBase(path string) {
    err := os.Remove(path)
    if err != nil {
        log.Fatal(err)
    }
}
