package handler

import "net/http"

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, apiError{
		Code:    code,
		Message: message,
	})
}

// exemplo de erro genérico (não vazar detalhe)
func writeInternalError(w http.ResponseWriter) {
	writeError(w, http.StatusInternalServerError, "internal_error", "something went wrong")
}
