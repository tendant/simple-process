package http

import (
	"encoding/json"
	"net/http"

	"github.com/tendant/simple-process/core/contracts"
)

// ResultHandler is an interface for handling UoW results.
type ResultHandler interface {
	HandleResult(http.ResponseWriter, *http.Request)
}

// DefaultResultHandler is a default implementation of ResultHandler.
type DefaultResultHandler struct{}

// HandleResult decodes the result from the request body and prints it.
// In a real application, this would trigger further processing.
func (h *DefaultResultHandler) HandleResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var result contracts.Result
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: Process the result (e.g., update metadata, trigger next UoW)
	// For now, just print it to the console.
	// log.Printf("Received result: %+v", result)

	w.WriteHeader(http.StatusOK)
}

