package database

import (
    "errors"
)

type Chirp struct {
    ID      int     `json:"id"`
    Body    string  `json:"body"`
    AuthorID int `json:"author_id"`
}

func (db *DB) DeleteChirp(chirpId int) error {
    dat, err := db.loadDB()
    if err != nil {
        return err
    }

    // get chirp
    _, err = db.GetChirp(chirpId)
    if err != nil {
        return err
    }

    delete(dat.Chirps, chirpId)
    
    return db.writeDB(dat)
}

func (db *DB) GetChirp(chirpID int) (Chirp, error) {
    dat, err := db.loadDB()
    if err != nil {
        return Chirp{}, err
    }

    chirp, ok := dat.Chirps[chirpID]
    if !ok {
        return Chirp{}, errors.New("No User Found") 
    }

    return chirp, nil
}

func (db *DB) CreateChirp(authorID int, body string) (Chirp, error) {
    chirp := Chirp{}
    dat, err := db.loadDB()
    if err != nil {
        return Chirp{}, err
    }

    id := len(dat.Chirps) + 1
    chirp.Body = body
    chirp.ID = id
    chirp.AuthorID = authorID
    dat.Chirps[id] = chirp

    db.writeDB(dat)
    return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
    dat, err := db.loadDB()
    if err != nil {
        return []Chirp{}, err 
    }

    chirps := make([]Chirp, 0, len(dat.Users))
    for _, v := range dat.Chirps {
        chirps = append(chirps, v)
    }

    return chirps, nil
}

