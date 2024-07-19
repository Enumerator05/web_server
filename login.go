package main

import (
    "net/http"
    "encoding/json"
    "web_server/internal/auth"
    "time"
)

func (apiConfig *ApiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
    type parameter struct {
        Password    string  `json:"password"`
        Email       string  `json:"email"`
        ExpiresIn   int     `json:"expires_in_seconds"`
    }

    type response struct {
        ID              int    `json:"id"`
        Email           string `json:"email"`
        Token           string `json:"token"`
        RefreshToken    string `json:"refresh_token"`
        IsChirpyRed     bool    `json:"is_chirpy_red"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameter{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
        return
    }

    user, err := apiConfig.db.GetUserByEmail(params.Email)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't find user")
        return
    }

    err = auth.ComparePassword(params.Password, user.Password)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Invalid Password")
        return
    }

    defaultExpiration := 60 * 60 * 24
    if params.ExpiresIn == 0 {
		params.ExpiresIn= defaultExpiration
	} else if params.ExpiresIn > defaultExpiration {
		params.ExpiresIn = defaultExpiration
    }

    token, err := auth.MakeJWT(user.ID, apiConfig.jwtSecret, time.Duration(params.ExpiresIn)*time.Second)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Could Not Make JWT")
        return
    }

    refreshToken, err := auth.MakeRefreshToken()
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Could not make refresh token")
        return
    }

    err = apiConfig.db.StoreRefreshToken(user.ID, refreshToken)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to store refresh token")
        return
    }

    respondWithJSON(w, http.StatusOK, 
        response{
            ID:             user.ID,
            Email:          user.Email,
            IsChirpyRed:      user.IsChirpyRed,
            Token:          token,
            RefreshToken:   refreshToken,
        })
}
