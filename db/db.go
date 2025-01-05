package db

import (
    "database/sql"
    "errors"
    "log"
    "time"
    _ "github.com/lib/pq"
)

var Connection *sql.DB

func InitDB() error {
    // Construct the connection string without fmt.Sprintf
    connStr := "host=localhost port=5432 user=postgres password=mysecretpassword dbname=postgres sslmode=disable"

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
            return errors.New("failed to initialize database: " + err.Error())
        }

        return nil
    }

    return errors.New("failed to connect to database after 5 attempts: " + err.Error())
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
        return errors.New("error creating drivers table: " + err.Error())
    }

    // Check if table is empty
    var count int
    err = Connection.QueryRow("SELECT COUNT(*) FROM drivers").Scan(&count)
    if err != nil {
        return errors.New("error checking drivers count: " + err.Error())
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
                return errors.New("error inserting test driver " + driver.Name + ": " + err.Error())
            }
        }
        log.Println("Successfully inserted test drivers into empty table")
    }

    return nil
}