package utils

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Data  interface{} `json:"data"`
	Error *string     `json:"error"`
}

func JSONResponse(w http.ResponseWriter, status int, data interface{}, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var errPtr *string
	if errMsg != "" {
		errPtr = &errMsg
	}

	resp := APIResponse{
		Data:  data,
		Error: errPtr,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to endcode data", http.StatusInternalServerError)
		return
	}
}
