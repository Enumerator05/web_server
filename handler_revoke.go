package main

import (
    "net/http"
    "web_server/internal/auth"
    "strconv"
)
func (cfg *ApiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
    tokenString, err := auth.GetBearerToken(r.Header)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Malformed headler")
        return 
    }

    userIdString, err := cfg.db.GetUserByRefreshToken(tokenString)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Could not validate refresh token")
        return
    } 

    userId, err := strconv.Atoi(userIdString)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Could not convert userId")
        return
    }

    user, err := cfg.db.GetUser(userId)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Could not get user")
        return
    }

    err = cfg.db.RevokeRefreshToken(&user)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Could not revoke refresh token")
        return
    }

    respondWithJSON(w, 204, nil)
}
