package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/health", handleHealth)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response["status"])
	}
}

func TestChangeEndpointValid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/change", handleChange)

	change := Change{
		Kind:       "Change",
		APIVersion: "v1",
		Spec: ChangeSpec{
			Prompt: "Add comprehensive error handling to all HTTP handlers",
			Repos:  []string{"https://github.com/myorg/repo1", "https://github.com/myorg/repo2"},
			Agent:  "copilot-cli",
			Branch: "main",
		},
	}

	jsonData, _ := json.Marshal(change)
	req, _ := http.NewRequest("POST", "/change", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "accepted" {
		t.Errorf("Expected status 'accepted', got '%v'", response["status"])
	}
}

func TestChangeEndpointDefaultBranch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/change", handleChange)

	change := Change{
		Kind:       "Change",
		APIVersion: "v1",
		Spec: ChangeSpec{
			Prompt: "Test prompt",
			Repos:  []string{"https://github.com/myorg/repo1"},
			Agent:  "gemini-cli",
			// Branch is omitted - should default to "main"
		},
	}

	jsonData, _ := json.Marshal(change)
	req, _ := http.NewRequest("POST", "/change", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	changeResponse := response["change"].(map[string]interface{})
	spec := changeResponse["spec"].(map[string]interface{})
	if spec["branch"] != "main" {
		t.Errorf("Expected default branch 'main', got '%v'", spec["branch"])
	}
}

func TestChangeEndpointInvalidKind(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/change", handleChange)

	invalidChange := map[string]interface{}{
		"kind":       "InvalidKind",
		"apiVersion": "v1",
		"spec": map[string]interface{}{
			"prompt": "Test",
			"repos":  []string{"https://github.com/myorg/repo1"},
			"agent":  "copilot-cli",
		},
	}

	jsonData, _ := json.Marshal(invalidChange)
	req, _ := http.NewRequest("POST", "/change", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestChangeEndpointInvalidAgent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/change", handleChange)

	change := Change{
		Kind:       "Change",
		APIVersion: "v1",
		Spec: ChangeSpec{
			Prompt: "Test",
			Repos:  []string{"https://github.com/myorg/repo1"},
			Agent:  "invalid-agent",
		},
	}

	jsonData, _ := json.Marshal(change)
	req, _ := http.NewRequest("POST", "/change", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Error != "invalid_agent" {
		t.Errorf("Expected error 'invalid_agent', got '%s'", response.Error)
	}
}

func TestChangeEndpointMissingPrompt(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/change", handleChange)

	invalidChange := map[string]interface{}{
		"kind":       "Change",
		"apiVersion": "v1",
		"spec": map[string]interface{}{
			"repos": []string{"https://github.com/myorg/repo1"},
			"agent": "copilot-cli",
		},
	}

	jsonData, _ := json.Marshal(invalidChange)
	req, _ := http.NewRequest("POST", "/change", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestChangeEndpointEmptyRepos(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/change", handleChange)

	change := Change{
		Kind:       "Change",
		APIVersion: "v1",
		Spec: ChangeSpec{
			Prompt: "Test",
			Repos:  []string{}, // Empty repos
			Agent:  "copilot-cli",
		},
	}

	jsonData, _ := json.Marshal(change)
	req, _ := http.NewRequest("POST", "/change", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Error != "missing_repos" {
		t.Errorf("Expected error 'missing_repos', got '%s'", response.Error)
	}
}
