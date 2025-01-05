package main

import (
    db "app/db"
    "encoding/json"
    "log"
    "net/http"
    "os"
    _ "path/filepath"
    "strings"

    "github.com/joho/godotenv"
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

type Ride struct {
    Id       string `json:"id"`
    CarId    string `json:"car_id"`
    Location string `json:"location"`
    Path     string `json:"path"`
}

func getRides(w http.ResponseWriter, req *http.Request) {
    rows, err := db.Connection.Query("SELECT * FROM rides")
    if err != nil {
        http.Error(w, "Failed to get rides: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var rides []Ride

    for rows.Next() {
        var ride Ride
        if err := rows.Scan(&ride.Id, &ride.CarId, &ride.Location, &ride.Path); err != nil {
            http.Error(w, "Error scanning ride: "+err.Error(), http.StatusInternalServerError)
            return
        }
        rides = append(rides, ride)
    }

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    
    ridesBytes, err := json.MarshalIndent(rides, "", "\t")
    if err != nil {
        http.Error(w, "Error encoding rides: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Write(ridesBytes)
}

func main() {
    // Get environment variables with defaults
    serverEnv := os.Getenv("SERVER_ENV")
    if serverEnv == "" {
        serverEnv = "DEV"
    }
    
    serverPort := os.Getenv("SERVER_PORT")
    if serverPort == "" {
        serverPort = "8080"
    }

    // Initialize database
    if err := db.InitDB(); err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer db.Connection.Close()

    // Create file server for static files
    fs := http.FileServer(http.Dir("frontend/build"))

    // Set up routes
    http.HandleFunc("/rides", getRides)
    http.HandleFunc("/drivers", getDrivers)

    // Handle all other routes
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Log the request for debugging
        log.Printf("Received request: %s %s", r.Method, r.URL.Path)

        // Set CORS headers
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

        // Handle OPTIONS requests
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        // Check if file exists
        path := "frontend/build" + r.URL.Path
        if _, err := os.Stat(path); os.IsNotExist(err) {
            // For API routes, don't serve index.html
            if strings.HasPrefix(r.URL.Path, "/drivers") || strings.HasPrefix(r.URL.Path, "/rides") {
                http.NotFound(w, r)
                return
            }

            // For all other routes, serve index.html
            log.Printf("File not found, serving index.html instead for path: %s", r.URL.Path)
            http.ServeFile(w, r, "frontend/build/index.html")
            return
        }

        // Set correct content types
        if strings.HasSuffix(r.URL.Path, ".js") {
            w.Header().Set("Content-Type", "application/javascript")
        } else if strings.HasSuffix(r.URL.Path, ".css") {
            w.Header().Set("Content-Type", "text/css")
        }

        // Serve static files
        fs.ServeHTTP(w, r)
    })

    // Log startup information
    log.Printf("Starting server in %s mode on port %s", serverEnv, serverPort)
    log.Printf("Serving static files from: %s", "frontend/build")

    // List contents of build directory
    if files, err := os.ReadDir("frontend/build"); err == nil {
        log.Println("Contents of build directory:")
        for _, file := range files {
            log.Printf("- %s", file.Name())
        }
    } else {
        log.Printf("Error reading build directory: %v", err)
    }

    // Start server with appropriate protocol
    var err error
    if serverEnv == "PROD" {
        // Use TLS in production
        err = http.ListenAndServeTLS(
            ":"+serverPort,
            "/etc/letsencrypt/live/app.evanomeje.xyz/fullchain.pem",
            "/etc/letsencrypt/live/app.evanomeje.xyz/privkey.pem",
            nil,
        )
    } else {
        // Use HTTP in development
        err = http.ListenAndServe(":"+serverPort, nil)
    }

    if err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
