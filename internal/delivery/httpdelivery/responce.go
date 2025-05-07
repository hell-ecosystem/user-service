package httpdelivery

import (
	"encoding/json"
	"net/http"
)

// APIResponse — общий формат ответа
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError — контейнер для ошибок
type APIError struct {
	Code    string `json:"code"`    // машинный код ошибки, e.g. "USER_NOT_FOUND"
	Message string `json:"message"` // человекочитаемое сообщение
}

// WriteJSON упрощает вывод любого APIResponse в JSON
func WriteJSON(w http.ResponseWriter, status int, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

// WriteSuccess возвращает { "success": true, "data": ... }
func WriteSuccess(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}

// WriteError возвращает { "success": false, "error": { "code": ..., "message": ... } }
func WriteError(w http.ResponseWriter, status int, code, message string) {
	WriteJSON(w, status, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}
