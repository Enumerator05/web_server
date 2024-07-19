package main

import (
    "net/http"
    "encoding/json"
    "web_server/internal/database"
    "web_server/internal/auth"
)
func (cfg *ApiConfig) handlerUpgradeToChirpyRed(w http.ResponseWriter, r *http.Request) {
    type Data struct {
        UserId  int     `json:"user_id"`
    }

    type Parameter struct {
        Event   string  `json:"event"`
        Data    Data    `json:"data"`
    }
    
    apiKey, err := auth.GetApiKey(r.Header)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Token not found")
        return
    }

    if apiKey != cfg.polkaApiKey {
        respondWithError(w, http.StatusUnauthorized, "Token not found")
        return
    }

    decoder := json.NewDecoder(r.Body)
    params := Parameter{}
    err = decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Failed to decode params")
        return
    }

    if params.Event != "user.upgraded" {
        respondWithError(w, 204, "Invalid Event")
        return
    }

    _, err = cfg.db.UpgradeUser(params.Data.UserId)
    if err != nil {
        respondWithError(w, http.StatusNotFound, "Could not upgrade user")
        return
    }

    respondWithJSON(w, 204, database.User{ })
}
