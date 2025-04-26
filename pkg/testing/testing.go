package testing

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

// TestTemplate is a basic template for Go tests.
const TestTemplate = `package {{.Package}}

import (
	"testing"
)

{{range .Functions}}
func Test{{.Name}}(t *testing.T) {
	{{if .TableDriven}}
	tests := []struct {
		name string
		// TODO: Add test case inputs and expected outputs
	}{
		{
			name: "test case 1",
		},
		{
			name: "test case 2",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Call {{.Name}} with the test case inputs and verify outputs
		})
	}
	{{else}}
	// TODO: Write test for {{.Name}}
	{{end}}
}
{{end}}
`

// TestData holds data for the test template.
type TestData struct {
	Package   string
	Functions []FunctionData
}

// FunctionData holds data about a function to test.
type FunctionData struct {
	Name        string
	TableDriven bool
}

// GenerateTests creates test files for Go functions.
func GenerateTests(path string, outputDir string, tableTests bool) error {
	fmt.Println("Generating tests for:", path)

	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if path is a directory
	fi, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("failed to stat path: %w", err)
	}

	if fi.IsDir() {
		// If it's a directory, process all Go files
		return filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
				return generateTestForFile(path, outputDir, tableTests)
			}

			return nil
		})
	} else if strings.HasSuffix(absPath, ".go") && !strings.HasSuffix(absPath, "_test.go") {
		// If it's a single Go file, process it
		return generateTestForFile(absPath, outputDir, tableTests)
	} else {
		return fmt.Errorf("path must be a directory or a Go file")
	}
}

// generateTestForFile creates a test file for a single Go file.
func generateTestForFile(path string, outputDir string, tableTests bool) error {
	// Parse the Go file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse Go file: %w", err)
	}

	// Get package name
	packageName := node.Name.Name

	// Find exported functions
	var functions []FunctionData
	for _, decl := range node.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && ast.IsExported(fn.Name.Name) {
			functions = append(functions, FunctionData{
				Name:        fn.Name.Name,
				TableDriven: tableTests,
			})
		}
	}

	if len(functions) == 0 {
		fmt.Printf("No exported functions found in %s, skipping\n", path)
		return nil
	}

	// Prepare output directory
	var outputPath string
	if outputDir == "" {
		// Use same directory as source file
		dir := filepath.Dir(path)
		baseName := filepath.Base(path)
		fileName := strings.TrimSuffix(baseName, ".go") + "_test.go"
		outputPath = filepath.Join(dir, fileName)
	} else {
		// Create output directory if it doesn't exist
		err = os.MkdirAll(outputDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		baseName := filepath.Base(path)
		fileName := strings.TrimSuffix(baseName, ".go") + "_test.go"
		outputPath = filepath.Join(outputDir, fileName)
	}

	// Check if test file already exists
	if _, err := os.Stat(outputPath); err == nil {
		return fmt.Errorf("test file already exists: %s", outputPath)
	}

	// Create template data
	data := TestData{
		Package:   packageName,
		Functions: functions,
	}

	// Parse and execute the template
	tmpl, err := template.New("test").Parse(TestTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse test template: %w", err)
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create test file: %w", err)
	}
	defer file.Close()

	// Execute the template
	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("failed to execute test template: %w", err)
	}

	fmt.Printf("Generated test file: %s\n", outputPath)
	return nil
}

// AnalyzeCoverage analyzes test coverage for a Go project.
func AnalyzeCoverage(path string, threshold float64, outputFile string) error {
	fmt.Printf("Analyzing test coverage for %s (threshold: %.1f%%)\n", path, threshold)

	// Get absolute paths
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	absOutput, err := filepath.Abs(outputFile)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for output: %w", err)
	}

	// Change to project directory
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(absPath)
	if err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}

	// Run tests with coverage
	coverProfilePath := "coverage.out"
	coverCmd := exec.Command("go", "test", "./...", "-coverprofile="+coverProfilePath)
	coverOutput, err := coverCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run tests with coverage: %w\nOutput: %s", err, coverOutput)
	}

	// Check if coverage file was created
	if _, err := os.Stat(coverProfilePath); os.IsNotExist(err) {
		return fmt.Errorf("coverage file was not created, ensure tests exist")
	}

	// Get coverage percentage
	funcCmd := exec.Command("go", "tool", "cover", "-func="+coverProfilePath)
	funcOutput, err := funcCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to analyze coverage: %w\nOutput: %s", err, funcOutput)
	}

	// Display coverage results
	fmt.Println("\nCoverage Results:")
	fmt.Println(string(funcOutput))

	// Generate HTML report
	htmlCmd := exec.Command("go", "tool", "cover", "-html="+coverProfilePath, "-o", absOutput)
	htmlOutput, err := htmlCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to generate HTML report: %w\nOutput: %s", err, htmlOutput)
	}

	// Extract total coverage percentage from output
	outputLines := strings.Split(string(funcOutput), "\n")
	var totalCoverage float64
	for _, line := range outputLines {
		if strings.Contains(line, "total:") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				percentStr := strings.TrimSuffix(parts[len(parts)-1], "%")
				fmt.Sscanf(percentStr, "%f", &totalCoverage)
				break
			}
		}
	}

	// Check if coverage meets threshold
	fmt.Printf("\nTotal coverage: %.1f%%\n", totalCoverage)
	fmt.Printf("Coverage HTML report generated at: %s\n", absOutput)

	if totalCoverage < threshold {
		fmt.Printf("\nWARNING: Coverage (%.1f%%) is below threshold (%.1f%%)\n", totalCoverage, threshold)
	} else {
		fmt.Printf("\nSUCCESS: Coverage (%.1f%%) meets or exceeds threshold (%.1f%%)\n", totalCoverage, threshold)
	}

	return nil
}
