package main

import (
    "net/http"
    "web_server/internal/auth"
    "time"
    "strconv"
)

func (apiConfig *ApiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
    type response struct {
        ID              int    `json:"id"`
        Email           string `json:"email"`
        Token           string `json:"token"`
        RefreshToken    string `json:"refresh_token"`
    }

    token, err := auth.GetBearerToken(r.Header)

    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Could not find refresh token")
        return
    }

    userIdString, err := apiConfig.db.GetUserByRefreshToken(token)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Could not validate refresh token")
        return
    } 

    userId, err := strconv.Atoi(userIdString)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Could not convert userId")
        return
    }

    user, err := apiConfig.db.GetUser(userId)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Could not get user")
        return
    }

    if user.RefreshToken.ExpiresAt.Before(time.Now().UTC()) {
        respondWithError(w, http.StatusUnauthorized, "Refresh token has expired")
        return
    }

    defaultExpiration := 60 * 60
    jwtToken, err := auth.MakeJWT(user.ID, apiConfig.jwtSecret, time.Duration(defaultExpiration)*time.Second)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Could Not Make JWT")
        return
    }

    respondWithJSON(w, http.StatusOK, 
        response{
            Token:          jwtToken,
        })
}
