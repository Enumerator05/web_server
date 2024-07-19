package main

import (
    "net/http"
    "fmt"
)

func (cfg *ApiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg.fileserverHits++
        next.ServeHTTP(w, r)
    })
}

func (cfg *ApiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(fmt.Sprintf(`<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>`, cfg.fileserverHits)))
    _ = r
}

func (cfg *ApiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    cfg.fileserverHits = 0
    _ = r
}
