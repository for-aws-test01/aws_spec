package common

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSuccessResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testData := map[string]string{"key": "value"}
	SuccessResponse(c, testData, "Success")

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response Response
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Message != "Success" {
		t.Errorf("Expected message 'Success', got '%s'", response.Message)
	}
}

func TestCreatedResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testData := map[string]string{"id": "123"}
	CreatedResponse(c, testData, "Created successfully")

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response Response
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Message != "Created successfully" {
		t.Errorf("Expected message 'Created successfully', got '%s'", response.Message)
	}
}

func TestPaginationResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	items := []string{"item1", "item2", "item3"}
	PaginationResponse(c, items, 10, 1, 3, "Success")

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response Response
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify pagination data structure
	dataMap, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Expected data to be a map")
	}

	if total, ok := dataMap["total"].(float64); !ok || int64(total) != 10 {
		t.Errorf("Expected total 10, got %v", dataMap["total"])
	}

	if page, ok := dataMap["page"].(float64); !ok || int(page) != 1 {
		t.Errorf("Expected page 1, got %v", dataMap["page"])
	}

	if totalPages, ok := dataMap["total_pages"].(float64); !ok || int(totalPages) != 4 {
		t.Errorf("Expected total_pages 4, got %v", dataMap["total_pages"])
	}
}

func TestErrorResponses(t *testing.T) {
	tests := []struct {
		name           string
		errorFunc      func(*gin.Context, string, string)
		expectedStatus int
	}{
		{"BadRequest", BadRequestError, http.StatusBadRequest},
		{"Unauthorized", UnauthorizedError, http.StatusUnauthorized},
		{"Forbidden", ForbiddenError, http.StatusForbidden},
		{"NotFound", NotFoundError, http.StatusNotFound},
		{"Conflict", ConflictError, http.StatusConflict},
		{"InternalServerError", InternalServerError, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.errorFunc(c, "Test error", "TEST_ERROR")

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			var response ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Error != "Test error" {
				t.Errorf("Expected error 'Test error', got '%s'", response.Error)
			}

			if response.Code != "TEST_ERROR" {
				t.Errorf("Expected code 'TEST_ERROR', got '%s'", response.Code)
			}
		})
	}
}
