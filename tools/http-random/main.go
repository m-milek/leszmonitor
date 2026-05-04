package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"time"

	"log/slog"
)

var possibleStatuses = []int{
	http.StatusOK,
	http.StatusBadRequest,
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		log.Info("request", "path", r.URL.Path, "status", 200)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Status: 200, Message: "healthy"})
	})

	mux.HandleFunc("GET /api/random", func(w http.ResponseWriter, r *http.Request) {
		code := possibleStatuses[rand.Intn(len(possibleStatuses))]
		log.Info("request", "path", r.URL.Path, "status", code)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(Response{Status: code, Message: http.StatusText(code)})
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Info("starting", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Error("fatal", "err", err)
		os.Exit(1)
	}
}
