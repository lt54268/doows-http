package api

import (
	"encoding/json"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		}
	}
}
