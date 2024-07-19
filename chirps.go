package main

import (
    "net/http"
    "log"
    "encoding/json"
    "fmt"
    "strings"
    "strconv"
    "web_server/internal/auth"
)

func (apiConfig *ApiConfig) POSTChirpsHandler(w http.ResponseWriter, r *http.Request) {
    type parameter struct {
        Body    string  `json:"body"`
    }

    tokenString, err := auth.GetBearerToken(r.Header)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "header has no token")
        return
    }

    userIdString, err := auth.ValidateJWT(tokenString, apiConfig.jwtSecret)
    if err != nil {
        userIdString, err = apiConfig.db.GetUserByRefreshToken(tokenString)

        if err != nil {
            respondWithError(w, http.StatusUnauthorized, "Could not validate token")
            return
        }
    }

    userId, err := strconv.Atoi(userIdString)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Could not parse userID")
        return
    }

    decoder := json.NewDecoder(r.Body)
    params := parameter{}
    err = decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Could not decode parameters")
        return
    }

    if len(params.Body) > 140 {
        respondWithError(w, http.StatusBadRequest, "Chirp is too long")
        return
    }
    
    type resp struct {
        Cleaned_body   string  `json:"cleaned_body"`
    }

    body := resp{Cleaned_body:  cleanText(params.Body)}
    chirp, err := apiConfig.db.CreateChirp(userId, body.Cleaned_body)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Could not create chirp")
        return
    }

    respondWithJSON(w, 201, chirp)
}

func (apiConfig *ApiConfig) GETChirpsHandler(w http.ResponseWriter, r *http.Request) {
    chirps, err := apiConfig.db.GetChirps()
    if err != nil {
        log.Fatal(err)
    }

    respondWithJSON(w, 200, chirps)
}

func (apiConfig *ApiConfig) HandlerDeleteChirpId(w http.ResponseWriter, r *http.Request) {
    token, err := auth.GetBearerToken(r.Header)   
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "No Bearer token found")
        return
    }

    chirpId, err := strconv.Atoi(r.PathValue("chirpID"))
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Could not parse chirpId to int")
        return
    }

    chirp, err := apiConfig.db.GetChirp(chirpId)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Could not fetch user")
        return
    }

    userIdString, err := auth.ValidateJWT(token, apiConfig.jwtSecret)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Could not validate token")
        return
    }

    userId, err := strconv.Atoi(userIdString)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Could not parse userIdString to int")
        return
    }

    if userId != chirp.ID {
        respondWithError(w, 403, "You do not own this chirp")
        return
    }

    err = apiConfig.db.DeleteChirp(chirpId)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Could not delete Chirp")
        return
    }

    respondWithJSON(w, 204, chirp)
}

func (apiConfig *ApiConfig) GETChirpsIDHandler(w http.ResponseWriter, r *http.Request) {
    chirps, err := apiConfig.db.GetChirps()
    if err != nil {
        log.Printf("error getting chirps")
        respondWithError(w, 500, "error getting chirps")
    }
    id, err := strconv.Atoi(r.PathValue("chirpID"))
    if err != nil {
        respondWithError(w, 404, "Invalid chirp ID")
        return
    }
    id--
    if id < 0 || id >= len(chirps) {
        respondWithError(w, 404, "Invalid chirp ID")
        return
    }

    chirp := chirps[id]
    respondWithJSON(w, 200, chirp)
}

func cleanText(txt string) string {
    badWords := [3]string{"kerfuffle", "sharbert", "fornax"}
    replaceTxt := "****"
    tokens := strings.Split(txt, " ")
    for i, token := range tokens {
        for _, badWord := range badWords {
            if badWord == strings.ToLower(token) {
                tokens[i] = replaceTxt
            }
        }
    }
    return strings.Join(tokens, " ")
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
    type error struct {
        Error   string  `json:"Error"`
    }
    errorResp := error{Error:   msg}
    respondWithJSON(w, code, errorResp)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    dat, err := json.Marshal(payload)

    if err != nil {
        w.WriteHeader(500)
        w.Write([]byte(fmt.Sprintf("Error marshalling JSON: %s", err)))
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write([]byte(dat))
}
