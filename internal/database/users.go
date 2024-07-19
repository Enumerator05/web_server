package database

import (
    "errors"
    "time"
    "fmt"
)

var (
   ErrUserNotFound = errors.New("User not found")
)

type User struct {
    ID              int     `json:"id"`
    Email           string  `json:"email"`
    Password        string  `json:"password"`
    RefreshToken    RefreshToken   `json:"refresh_token"`
    IsChirpyRed     bool    `json:"is_chirpy_red"`
}

type RefreshToken struct {
    TokenString       string `json:"token_string"`
    ExpiresAt   time.Time `json:"expires_at"`
}

func (db *DB) RevokeRefreshToken(user *User) error {
    dat, err := db.loadDB()
    if err != nil {
        return err
    }

    newUser := User {
        ID: user.ID,
        Email:  user.Email,
        Password:user.Password,
    }

    dat.Users[newUser.ID] = newUser
    return db.writeDB(dat)
}

func (db *DB) GetUserByRefreshToken(token string) (string, error) {
    dat, err := db.loadDB()
    if err != nil {
        return "", nil
    }

    for _, user := range dat.Users {
        if user.RefreshToken.TokenString == token {
            return fmt.Sprintf("%d", user.ID), nil
        }
    }

    return "", ErrUserNotFound
}

func (db *DB) StoreRefreshToken(userID int, tokenString string) error {
    // Load the database
    dat, err := db.loadDB()
    if err != nil {
        return err
    }

    // Create the refresh token
    token := RefreshToken{
        TokenString: tokenString,
        ExpiresAt:   time.Now().Add(60 * 24 * time.Hour).UTC(), // Expires in 60 days
    }

    // Retrieve the user
    user, ok := dat.Users[userID]
    if !ok {
        return ErrUserNotFound
    }

    // Assign the refresh token to the user
    user.RefreshToken = token

    // Update the user in the database
    dat.Users[userID] = user

    // Save the updated database
    err = db.writeDB(dat)
    if err != nil {
        return err
    }

    return nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
    dbstruct, err := db.loadDB()
    if err != nil {
        return User{}, err
    }
    
    for _, user := range dbstruct.Users {
        if user.Email == email {
            return user, nil
        }
    }

    return User{}, errors.New("User Not Found")
}

func (db *DB) UpgradeUser(userId int) (User, error) {
    dat, err := db.loadDB()

    if err != nil {
        return User{}, err
    }

    user, ok := dat.Users[userId]
    if !ok {
        return User{}, ErrUserNotFound
    }

    user.IsChirpyRed = true
    dat.Users[userId] = user

    err = db.writeDB(dat)
    if err != nil {
        return User{}, err
    }

    return user, nil
}

func (db *DB) UpdateEmailPassword(id int, email string, hashedPassword string) (User, error) {
    dat, err := db.loadDB()
    if err != nil {
        return User{}, err
    }

    user, ok := dat.Users[id]
    if !ok {
        return User{}, ErrUserNotFound
    }
    
    user = User {
        ID:         id,
        Email:      email,
        Password:   hashedPassword,
        RefreshToken: user.RefreshToken,
    }
    dat.Users[id] = user

    err = db.writeDB(dat)
    if err != nil {
        return User{}, err
    }

    return user, nil
}

func (db *DB) RevokeToken(id int) error {
    dat, err := db.loadDB()

    if err != nil {
        return err
    }

    user, ok := dat.Users[id]
    if !ok {
        return ErrUserNotFound
    }
    
    dat.Users[id] = User {
        ID:         user.ID,
        Email:      user.Email,
        Password:   user.Password,
    }

    err = db.writeDB(dat)
    if err != nil {
        return err
    }

    return nil
}

func (db *DB) CreateUser(email string, hashedPassword string) (User, error) {
    dat, err := db.loadDB()
    if err != nil {
        return User{}, err
    }

    _, err = db.GetUserByEmail(email)
    if err == nil {
        return User{}, errors.New("User Exists")
    }

    id := len(dat.Users) + 1

    user := User {
        ID:         id,
        Email:      email,
        Password:   hashedPassword,
        IsChirpyRed: false,
    }
    dat.Users[id] = user

    err = db.writeDB(dat)
    if err != nil {
        return User{}, err
    }

    return user, nil
}

func (db *DB) GetUsers() ([]User, error) {
    dat, err := db.loadDB()
    if err != nil {
        return []User{}, err 
    }

    users := make([]User, 0, len(dat.Users))
    for _, v := range dat.Users {
        users = append(users, v)
    }

    return users, nil
}

func (db *DB) GetUser(userID int) (User, error) {
    dat, err := db.loadDB()
    if err != nil {
        return User{}, err
    }

    user, ok := dat.Users[userID]
    if !ok {
        return User{}, errors.New("No User Found") 
    }

    return user, nil
}
