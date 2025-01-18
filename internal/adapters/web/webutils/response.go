package webutils

import (
	"encoding/json"
	"net/http"
	"time"
)

type Status string

const (
	Error   Status = "error"
	Success Status = "success"
)

type APIResponse struct {
	Status    Status      `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Metadata  interface{} `json:"metadata,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

func NewAPIResponse(
	status Status, message string,
	data interface{}, metadata interface{},
) APIResponse {

	return APIResponse{
		Status:    status,
		Message:   message,
		Data:      data,
		Metadata:  metadata,
		Timestamp: time.Now(),
	}
}
func WriteJSON(w http.ResponseWriter, status int, message string, data, metadata interface{}) error {
	var apiStatus Status
	if status >= 200 && status <= 299 {
		apiStatus = Success
	} else {
		apiStatus = Error
	}
	apiRes := NewAPIResponse(apiStatus, message, data, metadata)
	js, err := json.MarshalIndent(apiRes, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')
	// for key, value := range headers {
	// 	w.Header()[key] = value
	// }
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(js); err != nil {
		return err
	}
	return nil
}
