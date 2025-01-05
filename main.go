package main

import (
    db "app/postgres"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "path/filepath"
)

type Driver struct {
    ID            int    `json:"id"`
    Name          string `json:"name"`
    Phone         string `json:"phone"`
    Email         string `json:"email"`
    LicenseNumber string `json:"license_number"`
}

func init() {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
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

    // Create a file server handler
    fs := http.FileServer(http.Dir("./rideshare-frontend/build"))
    
    // Handle all routes
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Log the request
        log.Printf("Received request for: %s", r.URL.Path)

        // Check if the requested path exists
        requestedPath := filepath.Join("./rideshare-frontend/build", r.URL.Path)
        _, err := os.Stat(requestedPath)
        
        if err != nil {
            log.Printf("Error checking path %s: %v", requestedPath, err)
        }

        // If it's the root path or a non-existent file, serve index.html
        if r.URL.Path == "/" || os.IsNotExist(err) {
            indexPath := "./rideshare-frontend/build/index.html"
            log.Printf("Serving index.html from: %s", indexPath)
            
            // Check if index.html exists
            if _, err := os.Stat(indexPath); err != nil {
                log.Printf("Error: index.html not found at %s: %v", indexPath, err)
                http.Error(w, "index.html not found", http.StatusNotFound)
                return
            }
            
            http.ServeFile(w, r, indexPath)
            return
        }

        // For all other paths, try to serve the static file
        log.Printf("Serving static file: %s", requestedPath)
        fs.ServeHTTP(w, r)
    })

    // API routes
    http.HandleFunc("/drivers", getDrivers)

    serverPort := os.Getenv("SERVER_PORT")
    if serverPort == "" {
        serverPort = "8080"
    }
    serverEnv := os.Getenv("SERVER_ENV")

    log.Printf("Starting server in %s mode on port %s", serverEnv, serverPort)
    log.Printf("Static files directory: %s", "./rideshare-frontend/build")

    // List contents of build directory
    if files, err := os.ReadDir("./rideshare-frontend/build"); err == nil {
        log.Println("Contents of build directory:")
        for _, file := range files {
            log.Printf("- %s", file.Name())
        }
    } else {
        log.Printf("Error reading build directory: %v", err)
    }

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
