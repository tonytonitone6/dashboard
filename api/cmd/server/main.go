package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"dashboard/api/internal/auth"
	"dashboard/api/internal/server"
)

func main() {
	cfg, err := auth.ConfigFromEnv(os.Getenv)
	if err != nil {
		log.Fatal(err)
	}
	// make deploy runs this on the VM before it replaces the running binary.
	if len(os.Args) > 1 && os.Args[1] == "check-config" {
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("api listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, server.New(cfg, time.Now)))
}
