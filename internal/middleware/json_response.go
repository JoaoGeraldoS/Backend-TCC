package middleware

import (
	"encoding/json"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, status int, resul any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := json.NewEncoder(w).Encode(resul)
	return response
}
