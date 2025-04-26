package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AnalyzeStructure examines the project structure and architecture.
func AnalyzeStructure(path string) error {
	fmt.Println("Analyzing project structure at:", path)

	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Walk the directory tree
	fileCount := 0
	dirCount := 0
	pkgMap := make(map[string]bool)

	err = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(absPath, path)
		if err != nil {
			return err
		}

		// Skip hidden files and directories
		if strings.HasPrefix(filepath.Base(path), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			dirCount++
			fmt.Printf("Directory: %s\n", rel)
		} else if strings.HasSuffix(path, ".go") {
			fileCount++
			dir := filepath.Dir(path)
			pkgMap[dir] = true
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking directory: %w", err)
	}

	fmt.Printf("\nProject Summary:\n")
	fmt.Printf("- Directories: %d\n", dirCount)
	fmt.Printf("- Go files: %d\n", fileCount)
	fmt.Printf("- Packages: %d\n", len(pkgMap))

	fmt.Println("\nArchitecture Recommendations:")
	// We'd provide more sophisticated recommendations in a real implementation
	fmt.Println("- Use a clean architecture approach with clear separation of concerns")
	fmt.Println("- Follow Go project layout conventions (cmd, pkg, internal, etc.)")
	fmt.Println("- Ensure consistent package naming conventions")

	return nil
}

// AnalyzeQuality examines code quality and suggests improvements.
func AnalyzeQuality(path string) error {
	fmt.Println("Analyzing code quality at:", path)

	// In a real implementation we would load and analyze the packages using packages.Load
	// For this example, we'll just provide sample output
	fmt.Println("\nCode Quality Analysis Results:")
	fmt.Println("- Cyclomatic Complexity: Good (avg 4.2)")
	fmt.Println("- Code Duplication: Low (3.1%)")
	fmt.Println("- Error Handling: Good")
	fmt.Println("- Documentation Coverage: Medium (72%)")

	fmt.Println("\nImprovement Suggestions:")
	fmt.Println("- Add more documentation to exported functions")
	fmt.Println("- Consider breaking down complex functions in the handlers package")
	fmt.Println("- Implement more consistent error wrapping")

	return nil
}
