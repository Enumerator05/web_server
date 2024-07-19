package database

import (
    "sync"
    "encoding/json"
    "os"
    "log"
)

type DB struct {
    path    string  
    mux     *sync.RWMutex
}

type DBStructure struct {
    Users  map[int]User `json:"user"`
    Chirps  map[int]Chirp `json:"chips"`
}

func NewDB(path string) (*DB, error) {
    db := &DB{
        path:   path,
        mux:    &sync.RWMutex{},
        
    }
    err := db.ensureDB()
    if err != nil {
        return db, err
    }
    return db, nil
}

func (db *DB) ensureDB() error {
    _, err := os.Stat(db.path)
    if os.IsNotExist(err) {
        db.createDB()
        return nil
    }
    return err 
}

func (db *DB) createDB() {
    dbstructure := DBStructure{
        Users: map[int]User{},
        Chirps: map[int]Chirp{},
    }

    db.writeDB(dbstructure)
}

func (db *DB) loadDB() (DBStructure, error) {
    db.mux.RLock()
    defer db.mux.RUnlock()
    dat, err := os.ReadFile(db.path)
    dbstructure := DBStructure{}
    if err != nil {
        return DBStructure{}, err
    }

    err = json.Unmarshal(dat, &dbstructure)
    if err != nil {
        return DBStructure{}, err
    }

    return dbstructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
    db.mux.Lock()
    defer db.mux.Unlock()
    dat, err := json.Marshal(dbStructure)
    if err != nil {
        log.Fatal("Error marshallig json")
    }
    err = os.WriteFile(db.path, dat, 0600)
    if err != nil {
        return err
    }

    return nil
}
