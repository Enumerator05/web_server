package auth

import (
    "crypto/rand"
    "errors"
    "golang.org/x/crypto/bcrypt"
    "time"
    "github.com/golang-jwt/jwt/v5"
    "fmt"
    "net/http"
    "strings"
    "encoding/hex"
)

var (
   ErrNoAuthHeaderIncluded = errors.New("No auth header included in request") 
   ErrMalformedHeader = errors.New("malformed header") 
)

func HashPassword(password string) (string, error){
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

func ComparePassword(password, hashedPassword string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func MakeRefreshToken() (string, error) {
    size := 32
    b := make([]byte, size)

    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }

    token := hex.EncodeToString(b)
    return token, nil 
}

func MakeJWT(userID int, tokenSecret string, expiresIn time.Duration) (string, error) {
    signingKey := []byte(tokenSecret)

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
        Issuer: "chirpy",
        IssuedAt:   jwt.NewNumericDate(time.Now().UTC()),
        ExpiresAt:  jwt.NewNumericDate(time.Now().Add(expiresIn).UTC()),
        Subject:    fmt.Sprintf("%d", userID),
    })

    return token.SignedString(signingKey)
}

func ValidateJWT(tokenString, tokenSecret string) (string, error) {
    claims := jwt.RegisteredClaims{}
    token, err := jwt.ParseWithClaims(
        tokenString,
        &claims,
        func(token *jwt.Token) (interface{}, error) {return []byte(tokenSecret), nil},
    )

    if err != nil {
        return "", err
    }

    userIDString, err := token.Claims.GetSubject()
    if err != nil {
        return "", err
    }

    issuer, err := token.Claims.GetIssuer()
    if err != nil {
        return "", nil
    }

    if issuer != string("chirpy") {
        return "", errors.New("invalid issuer")
    }
    
    return userIDString, nil
}

func GetApiKey(header http.Header) (string, error) {
    authHeader := header.Get("Authorization")
    if authHeader == "" {
        return "", ErrNoAuthHeaderIncluded
    }

    splitAuth := strings.Split(authHeader, " ")
    if len(splitAuth) < 2 || splitAuth[0] != "ApiKey"{
        return "", ErrMalformedHeader
    }

    return splitAuth[1], nil
}


func GetBearerToken(header http.Header) (string, error) {
    authHeader := header.Get("Authorization")
    if authHeader == "" {
        return "", ErrNoAuthHeaderIncluded
    }

    splitAuth := strings.Split(authHeader, " ")
    if len(splitAuth) < 2 || splitAuth[0] != "Bearer"{
        return "", ErrMalformedHeader
    }

    return splitAuth[1], nil
}
