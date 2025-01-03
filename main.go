package main

import (
  "fmt"
  "net/http"
  "os"
  //db "rideshare-simulation/postgres"
  db "rideshare-simulation/postgres"
  //db "app/postgres"
)

/*
func getData(w http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(w, "Hello world\n")
}
*/

func getDrivers(w http.ResponseWriter, req *http.Request) {
	rows, err := db.Connection.Query("SELECT name FROM drivers")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	data := ""
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(name)
		data += fmt.Sprintf("%s ", name)
	}

	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Fprintf(w, data)
}

func main() {
  db.InitDB()
  defer db.Connection.Close()

  //http.HandleFunc("/data", getData)
http.Handle("/", http.FileServer(http.Dir("./static")))
http.HandleFunc("/drivers", getDrivers)

  // Get the SERVER_ENV environment variable
  serverEnv := os.Getenv("SERVER_ENV")

  // Run the server based on the environment
  if serverEnv == "DEV" {
    fmt.Println("Running in DEV mode on port 8080...")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
      fmt.Println("Error starting HTTP server:", err)
    }
  } else if serverEnv == "PROD" {
    fmt.Println("Running in PROD mode on port 443...")
    err := http.ListenAndServeTLS(
      ":443",
      "/etc/letsencrypt/live/app.evanomeje.xyz/fullchain.pem",
      "/etc/letsencrypt/live/app.evanomeje.xyz/privkey.pem",
      nil,
    )
    if err != nil {
      fmt.Println("Error starting HTTPS server:", err)
    }
  } else {
    fmt.Println("Unknown SERVER_ENV value. Please set it to DEV or PROD.")
  }
}
