package main

import (
	"log"
	"net/http"
	"os"

	"dashboard/api/internal/overview"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})
	overview.Register(mux)

	log.Printf("api listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
