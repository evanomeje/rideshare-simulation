package main

import (
  "fmt"
  "net/http"
  "os"
)

func getData(w http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(w, "Hello world\n")
}

func main() {
  http.HandleFunc("/data", getData)
http.Handle("/", http.FileServer(http.Dir("./static")))

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
