package httpx

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// ErrorResponse is the standard error envelope.
type ErrorResponse struct {
	Error     ErrorDetail `json:"error"`
	RequestID string      `json:"request_id,omitempty"`
}

// ErrorDetail contains error code and message.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Response is a generic success response envelope.
type Response struct {
	Data      interface{} `json:"data,omitempty"`
	Meta      interface{} `json:"meta,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

// WriteJSON writes a JSON response with the given status code.
func WriteJSON(w http.ResponseWriter, r *http.Request, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			// If encoding fails, try to write a minimal error
			http.Error(w, `{"error":{"code":"internal_error","message":"failed to encode response"}}`, http.StatusInternalServerError)
		}
	}
}

// WriteData writes a success response with data wrapped in the standard envelope.
func WriteData(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	requestID := middleware.GetReqID(r.Context())

	resp := Response{
		Data:      data,
		RequestID: requestID,
	}

	WriteJSON(w, r, status, resp)
}

// WriteDataWithMeta writes a success response with data and metadata.
func WriteDataWithMeta(w http.ResponseWriter, r *http.Request, status int, data interface{}, meta interface{}) {
	requestID := middleware.GetReqID(r.Context())

	resp := Response{
		Data:      data,
		Meta:      meta,
		RequestID: requestID,
	}

	WriteJSON(w, r, status, resp)
}

// WriteError writes an error response with the standard envelope.
func WriteError(w http.ResponseWriter, r *http.Request, status int, code string, message string) {
	requestID := middleware.GetReqID(r.Context())

	resp := ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
		RequestID: requestID,
	}

	WriteJSON(w, r, status, resp)
}

// Common error helpers

// WriteBadRequest writes a 400 Bad Request error.
func WriteBadRequest(w http.ResponseWriter, r *http.Request, message string) {
	WriteError(w, r, http.StatusBadRequest, "bad_request", message)
}

// WriteUnauthorized writes a 401 Unauthorized error.
func WriteUnauthorized(w http.ResponseWriter, r *http.Request, message string) {
	WriteError(w, r, http.StatusUnauthorized, "unauthorized", message)
}

// WriteForbidden writes a 403 Forbidden error.
func WriteForbidden(w http.ResponseWriter, r *http.Request, message string) {
	WriteError(w, r, http.StatusForbidden, "forbidden", message)
}

// WriteNotFound writes a 404 Not Found error.
func WriteNotFound(w http.ResponseWriter, r *http.Request, message string) {
	WriteError(w, r, http.StatusNotFound, "not_found", message)
}

// WriteConflict writes a 409 Conflict error.
func WriteConflict(w http.ResponseWriter, r *http.Request, message string) {
	WriteError(w, r, http.StatusConflict, "conflict", message)
}

// WriteUnprocessable writes a 422 Unprocessable Entity error.
func WriteUnprocessable(w http.ResponseWriter, r *http.Request, message string) {
	WriteError(w, r, http.StatusUnprocessableEntity, "unprocessable_entity", message)
}

// WriteInternalError writes a 500 Internal Server Error.
func WriteInternalError(w http.ResponseWriter, r *http.Request, message string) {
	WriteError(w, r, http.StatusInternalServerError, "internal_error", message)
}

// WriteServiceUnavailable writes a 503 Service Unavailable error.
func WriteServiceUnavailable(w http.ResponseWriter, r *http.Request, message string) {
	WriteError(w, r, http.StatusServiceUnavailable, "service_unavailable", message)
}

// NoContent writes a 204 No Content response.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

