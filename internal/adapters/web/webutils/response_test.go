package webutils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestWriteJSON_Success(t *testing.T) {
	rr := httptest.NewRecorder()

	data := map[string]string{"key": "value"}
	metadata := map[string]string{"metaKey": "metaValue"}
	err := WriteJSON(rr, http.StatusOK, "Success message", data, metadata)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, status)
	}

	var actualResponse APIResponse
	err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// Normalize and compare
	if !reflect.DeepEqual(toInterfaceMap(data), actualResponse.Data) {
		t.Fatalf("Data mismatch:\nExpected: %+v\nActual: %+v", data, actualResponse.Data)
	}
	if actualResponse.Message != "Success message" {
		t.Fatalf("Message mismatch: expected 'Success message', got %s", actualResponse.Message)
	}
	if actualResponse.Status != Success {
		t.Fatalf("Status mismatch: expected 'Success', got %s", actualResponse.Status)
	}
	if !reflect.DeepEqual(toInterfaceMap(metadata), actualResponse.Metadata) {
		t.Fatalf("Metadata mismatch:\nExpected: %+v\nActual: %+v", metadata, actualResponse.Metadata)
	}

	// Validate Timestamp
	if time.Since(actualResponse.Timestamp) > time.Second {
		t.Fatalf("Timestamp mismatch: expected recent, got %v", actualResponse.Timestamp)
	}
}

func toInterfaceMap(input map[string]string) map[string]interface{} {
	normalized := make(map[string]interface{}, len(input))
	for key, value := range input {
		normalized[key] = value
	}
	return normalized
}

func TestWriteJSON_Error(t *testing.T) {

	rr := httptest.NewRecorder()

	err := WriteJSON(rr, http.StatusInternalServerError, "Error message", nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Fatalf("expected status code %d, got %d", http.StatusInternalServerError, status)
	}

	expectedResponse := APIResponse{
		Status:    Error,
		Message:   "Error message",
		Data:      nil,
		Metadata:  nil,
		Timestamp: time.Now(),
	}

	var actualResponse APIResponse
	err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if actualResponse.Status != expectedResponse.Status ||
		actualResponse.Message != expectedResponse.Message ||
		actualResponse.Data != expectedResponse.Data ||
		actualResponse.Metadata != expectedResponse.Metadata {
		t.Fatalf("expected %+v, got %+v", expectedResponse, actualResponse)
	}

	if time.Since(actualResponse.Timestamp) > time.Second {
		t.Fatalf("expected timestamp to be recent, got %v", actualResponse.Timestamp)
	}
}
