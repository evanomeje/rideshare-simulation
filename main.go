package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    db "app/postgres"
)

type Driver struct {
    ID            int    `json:"id"`
    Name          string `json:"name"`
    Phone         string `json:"phone"`
    Email         string `json:"email"`
    LicenseNumber string `json:"license_number"`
}

func getDrivers(w http.ResponseWriter, req *http.Request) {
    rows, err := db.Connection.Query(`
        SELECT id, name, phone, email, license_number 
        FROM drivers
    `)
    if err != nil {
        log.Printf("Database query error: %v", err)
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var drivers []Driver
    for rows.Next() {
        var d Driver
        if err := rows.Scan(&d.ID, &d.Name, &d.Phone, &d.Email, &d.LicenseNumber); err != nil {
            log.Printf("Row scan error: %v", err)
            http.Error(w, "Database error", http.StatusInternalServerError)
            return
        }
        drivers = append(drivers, d)
    }

    if err = rows.Err(); err != nil {
        log.Printf("Row iteration error: %v", err)
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }

    // Set response headers
    w.Header().Set("Content-Type", "application/json")
    
    // Encode and send the response
    if err := json.NewEncoder(w).Encode(drivers); err != nil {
        log.Printf("Error encoding response: %v", err)
        http.Error(w, "Error encoding response", http.StatusInternalServerError)
        return
    }
}

func main() {
    if err := db.InitDB(); err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer db.Connection.Close()

    http.Handle("/", http.FileServer(http.Dir("./static")))
    http.HandleFunc("/drivers", getDrivers)

    serverPort := os.Getenv("SERVER_PORT")
    if serverPort == "" {
        serverPort = "8080"
    }
    serverEnv := os.Getenv("SERVER_ENV")

    log.Printf("Starting server in %s mode on port %s", serverEnv, serverPort)
    
    var err error
    if serverEnv == "PROD" {
        err = http.ListenAndServeTLS(
            ":"+serverPort,
            "/etc/letsencrypt/live/app.evanomeje.xyz/fullchain.pem",
            "/etc/letsencrypt/live/app.evanomeje.xyz/privkey.pem",
            nil,
        )
    } else {
        err = http.ListenAndServe(":"+serverPort, nil)
    }

    if err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
