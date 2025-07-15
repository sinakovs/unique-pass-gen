package main

import (
	"log"
	"net/http"
	"time"

	"unique-pass-gen/internal/router"
	"unique-pass-gen/internal/storage"
)

var port = "8080"

func main() {
	passwordCache := storage.NewCache()

	handler := router.Routes(passwordCache)

	srv := &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
