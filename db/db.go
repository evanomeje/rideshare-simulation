package db

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"
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
    configFile, err := os.Open("../dbconfig.json")
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

    for i := 0; i < 5; i++ {
        log.Printf("Attempting database connection (attempt %d/5)...", i+1)
        Connection, err = sql.Open("postgres", connStr)
        if err != nil {
            log.Printf("Error opening database: %v", err)
            time.Sleep(time.Second * 5)
            continue
        }

        err = Connection.Ping()
        if err != nil {
            log.Printf("Error pinging database: %v", err)
            time.Sleep(time.Second * 5)
            continue
        }

        log.Println("Successfully connected to the database!")
        return nil
    }

    return fmt.Errorf("failed to connect to database after 5 attempts: %v", err)
}