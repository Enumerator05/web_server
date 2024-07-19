package main

import (
    "net/http"
    "encoding/json"
    "web_server/internal/auth"
    "strconv"
    "web_server/internal/database"
    "log"
)

func (apiConfig *ApiConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    type parameter struct {
        Email       string  `json:"email"`
        Password    string  `json:"password"`
        ExpiresIn   int     `json:"exires_in_seconds"`
    }

    type response struct {
        User    database.User    `json:"user"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameter{}
    err := decoder.Decode(&params)
    if err != nil {
        w.WriteHeader(500)
    }

    hashedPassword, err := auth.HashPassword(params.Password)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Could not hash password")
    }

    user, err := apiConfig.db.CreateUser(params.Email, hashedPassword)
    if err != nil {
        log.Print(err)
        respondWithError(w, http.StatusNotAcceptable, "User Not Created")
        return
    }

    respondWithJSON(w, 201, database.User{
        ID: user.ID,
        Email:  user.Email,
        IsChirpyRed: user.IsChirpyRed,
    })
}

func (apiConfig *ApiConfig) handlerUpdateUserInformation(w http.ResponseWriter, r *http.Request) {
    type parameter struct {
        Email    string  `json:"email"`
        Password string  `json:"password"`
    }

    token, err := auth.GetBearerToken(r.Header)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
        return
    }
    
    subject, err := auth.ValidateJWT(token, apiConfig.jwtSecret)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
        return
    }

    decoder := json.NewDecoder(r.Body)
    params := parameter{}
    err = decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Couldn't decode parameters")
        return
    }

    hashedPassword, err := auth.HashPassword(params.Password)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "couldn't hash password")
        return
    }

    userIDInt, err := strconv.Atoi(subject)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
        return
    }

    _, err = apiConfig.db.GetUser(userIDInt)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't find user")
        return
    }

    user, err := apiConfig.db.UpdateEmailPassword(userIDInt, params.Email, hashedPassword)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
        return
    }

    respondWithJSON(w, http.StatusOK, database.User{
        ID:     user.ID,
        Email:  user.Email,
        IsChirpyRed: user.IsChirpyRed,
    })
}
