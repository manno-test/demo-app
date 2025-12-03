package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// ChangeSpec defines the specification for a change request
type ChangeSpec struct {
	Prompt string   `json:"prompt" binding:"required"`
	Repos  []string `json:"repos" binding:"required"`
	Agent  string   `json:"agent" binding:"required"`
	Branch string   `json:"branch"`
}

// Change represents the entire change request
type Change struct {
	Kind       string     `json:"kind" binding:"required,eq=Change"`
	APIVersion string     `json:"apiVersion" binding:"required"`
	Spec       ChangeSpec `json:"spec" binding:"required"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

var logger *slog.Logger

func init() {
	// Initialize slog logger with JSON handler
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

func main() {
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Add custom middleware for logging and recovery
	router.Use(ginLogger(), gin.Recovery())

	// Register routes
	router.POST("/change", handleChange)
	router.GET("/health", handleHealth)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Starting API server", "port", port)

	if err := router.Run(":" + port); err != nil {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}

// ginLogger is a middleware that logs requests using slog
func ginLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Log after processing
		statusCode := c.Writer.Status()
		logger.Info("Request processed",
			"method", method,
			"path", path,
			"status", statusCode,
			"ip", c.ClientIP(),
		)
	}
}

// handleHealth handles health check requests
func handleHealth(c *gin.Context) {
	logger.Info("Health check requested")

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "demo-app",
	})
}

// handleChange handles change request submissions
func handleChange(c *gin.Context) {
	var change Change

	// Bind and validate JSON
	if err := c.ShouldBindJSON(&change); err != nil {
		logger.Error("Failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Validate kind field
	if change.Kind != "Change" {
		logger.Warn("Invalid kind field", "kind", change.Kind)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_kind",
			Message: "kind must be 'Change'",
		})
		return
	}

	// Validate API version
	if change.APIVersion == "" {
		logger.Warn("Missing apiVersion field")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "missing_api_version",
			Message: "apiVersion is required",
		})
		return
	}

	// Validate spec fields
	if change.Spec.Prompt == "" {
		logger.Warn("Missing prompt in spec")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "missing_prompt",
			Message: "spec.prompt is required",
		})
		return
	}

	if len(change.Spec.Repos) == 0 {
		logger.Warn("No repositories specified")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "missing_repos",
			Message: "spec.repos must contain at least one repository",
		})
		return
	}

	if change.Spec.Agent == "" {
		logger.Warn("Missing agent in spec")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "missing_agent",
			Message: "spec.agent is required",
		})
		return
	}

	// Validate agent value
	validAgents := map[string]bool{
		"copilot-cli": true,
		"gemini-cli":  true,
	}
	if !validAgents[change.Spec.Agent] {
		logger.Warn("Invalid agent specified", "agent", change.Spec.Agent)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_agent",
			Message: "spec.agent must be either 'copilot-cli' or 'gemini-cli'",
		})
		return
	}

	// Set default branch if not provided
	if change.Spec.Branch == "" {
		change.Spec.Branch = "main"
		logger.Info("Using default branch", "branch", "main")
	}

	// Log successful change request
	logger.Info("Change request received",
		"prompt", change.Spec.Prompt,
		"repos", change.Spec.Repos,
		"agent", change.Spec.Agent,
		"branch", change.Spec.Branch,
	)

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"status":  "accepted",
		"message": "Change request received successfully",
		"change":  change,
	})
}
