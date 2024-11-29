package utils

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(ctx MyContext, w http.ResponseWriter, message string) {
	ctx.Logger.Error(message)

	var statusCode int
	if strings.Contains(message, "not found") {
		statusCode = http.StatusNotFound
	} else if strings.Contains(message, "missing") || strings.Contains(message, "invalid") {
		statusCode = http.StatusBadRequest
	} else {
		statusCode = http.StatusInternalServerError
	}

	errRes := ErrorResponse{Message: message}

	jsonErrRes, err := json.Marshal(errRes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(ContentType, ApplicationJSON)
	w.WriteHeader(statusCode)
	w.Write(jsonErrRes)
}
