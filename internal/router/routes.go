package router

import (
	"net/http"

	"unique-pass-gen/internal/handler"
	"unique-pass-gen/internal/storage"
)

func Routes(cache *storage.Cache) http.Handler {
	api := http.NewServeMux()
	h := handler.NewHandler(cache)

	api.HandleFunc("GET /", h.GetForm)
	api.HandleFunc("POST /", h.GeneratePass)

	return api
}
