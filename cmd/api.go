package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"goforge/pkg/analyzer"
	"goforge/pkg/dependency"
	"goforge/pkg/docs"

	"github.com/urfave/cli/v2"
)

// APICommand returns the CLI command for starting the API server.
func APICommand() *cli.Command {
	return &cli.Command{
		Name:  "api",
		Usage: "Start the API server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   "8080",
				Usage:   "Port to run the API server on",
			},
		},
		Action: func(c *cli.Context) error {
			port := c.String("port")
			return startAPIServer(port)
		},
	}
}

// ErrorResponse represents an error response from the API.
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a success response from the API.
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// startAPIServer starts the API server on the specified port.
func startAPIServer(port string) error {
	fmt.Printf("Starting API server on port %s...\n", port)

	// Define API routes
	http.HandleFunc("/api/health", healthCheckHandler)
	http.HandleFunc("/api/analyze/structure", analyzeStructureHandler)
	http.HandleFunc("/api/analyze/quality", analyzeQualityHandler)
	http.HandleFunc("/api/dependency/check", checkDependenciesHandler)
	http.HandleFunc("/api/docs/generate", generateDocsHandler)

	// Start the server
	addr := ":" + port
	fmt.Printf("API server is running at http://localhost%s\n", addr)
	fmt.Println("Press Ctrl+C to stop")
	return http.ListenAndServe(addr, nil)
}

// healthCheckHandler handles health check requests.
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := SuccessResponse{
		Message: "GoForge API is running",
		Data: map[string]string{
			"status":  "healthy",
			"version": "1.0.0",
		},
	}

	sendJSON(w, response, http.StatusOK)
}

// analyzeStructureHandler handles requests to analyze project structure.
func analyzeStructureHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request
	err := r.ParseForm()
	if err != nil {
		sendError(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	path := r.FormValue("path")
	if path == "" {
		sendError(w, "Path is required", http.StatusBadRequest)
		return
	}

	// Create a temporary file to capture output
	tempFile, err := os.CreateTemp("", "goforge-api-*.txt")
	if err != nil {
		sendError(w, "Failed to create temporary file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Redirect stdout to the temporary file
	oldStdout := os.Stdout
	os.Stdout = tempFile
	defer func() { os.Stdout = oldStdout }()

	// Run the analysis
	err = analyzer.AnalyzeStructure(path)
	if err != nil {
		sendError(w, fmt.Sprintf("Failed to analyze structure: %v", err), http.StatusInternalServerError)
		return
	}

	// Reset file pointer and read the output
	tempFile.Seek(0, 0)
	output, err := os.ReadFile(tempFile.Name())
	if err != nil {
		sendError(w, "Failed to read analysis output", http.StatusInternalServerError)
		return
	}

	// Send the response
	response := SuccessResponse{
		Message: "Project structure analyzed successfully",
		Data: map[string]string{
			"output": string(output),
		},
	}

	sendJSON(w, response, http.StatusOK)
}

// analyzeQualityHandler handles requests to analyze code quality.
func analyzeQualityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request
	err := r.ParseForm()
	if err != nil {
		sendError(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	path := r.FormValue("path")
	if path == "" {
		sendError(w, "Path is required", http.StatusBadRequest)
		return
	}

	// Create a temporary file to capture output
	tempFile, err := os.CreateTemp("", "goforge-api-*.txt")
	if err != nil {
		sendError(w, "Failed to create temporary file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Redirect stdout to the temporary file
	oldStdout := os.Stdout
	os.Stdout = tempFile
	defer func() { os.Stdout = oldStdout }()

	// Run the analysis
	err = analyzer.AnalyzeQuality(path)
	if err != nil {
		sendError(w, fmt.Sprintf("Failed to analyze quality: %v", err), http.StatusInternalServerError)
		return
	}

	// Reset file pointer and read the output
	tempFile.Seek(0, 0)
	output, err := os.ReadFile(tempFile.Name())
	if err != nil {
		sendError(w, "Failed to read analysis output", http.StatusInternalServerError)
		return
	}

	// Send the response
	response := SuccessResponse{
		Message: "Code quality analyzed successfully",
		Data: map[string]string{
			"output": string(output),
		},
	}

	sendJSON(w, response, http.StatusOK)
}

// checkDependenciesHandler handles requests to check dependencies.
func checkDependenciesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request
	err := r.ParseForm()
	if err != nil {
		sendError(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	path := r.FormValue("path")
	if path == "" {
		sendError(w, "Path is required", http.StatusBadRequest)
		return
	}

	// Create a temporary file to capture output
	tempFile, err := os.CreateTemp("", "goforge-api-*.txt")
	if err != nil {
		sendError(w, "Failed to create temporary file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Redirect stdout to the temporary file
	oldStdout := os.Stdout
	os.Stdout = tempFile
	defer func() { os.Stdout = oldStdout }()

	// Run the dependency check
	err = dependency.CheckOutdated(path)
	if err != nil {
		sendError(w, fmt.Sprintf("Failed to check dependencies: %v", err), http.StatusInternalServerError)
		return
	}

	// Reset file pointer and read the output
	tempFile.Seek(0, 0)
	output, err := os.ReadFile(tempFile.Name())
	if err != nil {
		sendError(w, "Failed to read dependency check output", http.StatusInternalServerError)
		return
	}

	// Send the response
	response := SuccessResponse{
		Message: "Dependencies checked successfully",
		Data: map[string]string{
			"output": string(output),
		},
	}

	sendJSON(w, response, http.StatusOK)
}

// generateDocsHandler handles requests to generate documentation.
func generateDocsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request
	err := r.ParseForm()
	if err != nil {
		sendError(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	path := r.FormValue("path")
	if path == "" {
		sendError(w, "Path is required", http.StatusBadRequest)
		return
	}

	docType := r.FormValue("type")
	if docType == "" {
		docType = "user" // Default to user docs
	}

	format := r.FormValue("format")
	if format == "" {
		format = "markdown" // Default to markdown
	}

	outputDir := r.FormValue("output")
	if outputDir == "" {
		outputDir = filepath.Join(os.TempDir(), "goforge-docs")
	}

	// Create a temporary file to capture output
	tempFile, err := os.CreateTemp("", "goforge-api-*.txt")
	if err != nil {
		sendError(w, "Failed to create temporary file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Redirect stdout to the temporary file
	oldStdout := os.Stdout
	os.Stdout = tempFile
	defer func() { os.Stdout = oldStdout }()

	// Generate the documentation
	var docErr error
	if docType == "api" {
		docErr = docs.GenerateAPIDoc(path, outputDir, format)
	} else {
		docErr = docs.GenerateUserDoc(path, outputDir, format)
	}

	if docErr != nil {
		sendError(w, fmt.Sprintf("Failed to generate documentation: %v", docErr), http.StatusInternalServerError)
		return
	}

	// Reset file pointer and read the output
	tempFile.Seek(0, 0)
	output, err := os.ReadFile(tempFile.Name())
	if err != nil {
		sendError(w, "Failed to read documentation output", http.StatusInternalServerError)
		return
	}

	// Send the response
	response := SuccessResponse{
		Message: "Documentation generated successfully",
		Data: map[string]interface{}{
			"output":    string(output),
			"directory": outputDir,
		},
	}

	sendJSON(w, response, http.StatusOK)
}

// sendJSON sends a JSON response with the given status code.
func sendJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

// sendError sends an error response with the given status code.
func sendError(w http.ResponseWriter, message string, status int) {
	response := ErrorResponse{
		Error: message,
	}

	sendJSON(w, response, status)
}
