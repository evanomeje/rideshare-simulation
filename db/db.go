package db

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "os"
    _ "github.com/lib/pq"
)

type Config struct {
    User     string `json:"user"`
    Password string `json:"password"`
    Host     string `json:"host"`
    Port     int    `json:"port"`
    DBname   string `json:"dbname"`
}

var Connection *sql.DB

func InitDB() error {
    configFile, err := os.Open("dbconfig.json")
    if err != nil {
        return fmt.Errorf("error opening dbconfig.json: %v", err)
    }
    defer configFile.Close()

    var config Config
    jsonParser := json.NewDecoder(configFile)
    if err = jsonParser.Decode(&config); err != nil {
        return fmt.Errorf("error decoding dbconfig.json: %v", err)
    }

    connStr := fmt.Sprintf(
        "user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
        config.User, config.Password, config.Host, config.Port, config.DBname,
    )

    Connection, err = sql.Open("postgres", connStr)
    if err != nil {
        return fmt.Errorf("error connecting to database: %v", err)
    }

    return Connection.Ping()
}