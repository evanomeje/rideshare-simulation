package db

import (
    "database/sql"
    "fmt"
    "log"
    "time"
    _ "github.com/lib/pq"
)

var Connection *sql.DB

func InitDB() error {
    connStr := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        "db",           // host - matches service name in docker-compose
        5432,          // port
        "postgres",    // user
        "mysecretpassword", // password - should use env var in production
        "postgres",    // dbname
    )

    var err error
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
        
        if err := initializeDatabase(); err != nil {
            return fmt.Errorf("failed to initialize database: %v", err)
        }
        
        return nil
    }

    return fmt.Errorf("failed to connect to database after 5 attempts: %v", err)
}

type Driver struct {
    Name          string
    Phone         string
    Email         string
    Password      string
    LicenseNumber string
}

func initializeDatabase() error {
    // Create the drivers table
    createTableSQL := `
    CREATE TABLE IF NOT EXISTS drivers (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        phone VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL,
        password VARCHAR(255) NOT NULL,
        license_number VARCHAR(255) NOT NULL
    );`

    _, err := Connection.Exec(createTableSQL)
    if err != nil {
        return fmt.Errorf("error creating drivers table: %v", err)
    }

    // Check if table is empty
    var count int
    err = Connection.QueryRow("SELECT COUNT(*) FROM drivers").Scan(&count)
    if err != nil {
        return fmt.Errorf("error checking drivers count: %v", err)
    }

    // Insert test data if table is empty
    if count == 0 {
        testDrivers := []Driver{
            {
                Name:          "Alice",
                Phone:         "+14031234567",
                Email:         "alice@example.com",
                Password:      "b7Fkd9Lm",
                LicenseNumber: "XX-ZZ-23",
            },
            {
                Name:          "Michael",
                Phone:         "+15873987654",
                Email:         "michael@example.com",
                Password:      "a5Sgf8Wx",
                LicenseNumber: "BB-CC-45",
            },
            {
                Name:          "Nancy",
                Phone:         "+14039876543",
                Email:         "nancy@example.com",
                Password:      "n3TcH1Ld",
                LicenseNumber: "YY-AA-67",
            },
            {
                Name:          "Sarah",
                Phone:         "+15873216587",
                Email:         "sarah@example.com",
                Password:      "x9KfY2Vr",
                LicenseNumber: "ZZ-DD-89",
            },
        }

        insertSQL := `
        INSERT INTO drivers (name, phone, email, password, license_number)
        VALUES ($1, $2, $3, $4, $5)`

        for _, driver := range testDrivers {
            _, err = Connection.Exec(insertSQL,
                driver.Name,
                driver.Phone,
                driver.Email,
                driver.Password,
                driver.LicenseNumber,
            )
            if err != nil {
                return fmt.Errorf("error inserting test driver %s: %v", driver.Name, err)
            }
        }
        log.Println("Successfully inserted test drivers into empty table")
    }

    return nil
}
