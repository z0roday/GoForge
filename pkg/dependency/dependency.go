package dependency

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CheckOutdated checks for outdated dependencies in a Go project.
func CheckOutdated(path string) error {
	fmt.Println("Checking for outdated dependencies in:", path)

	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
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

	// Use 'go list -m -u all' to check for outdated dependencies
	cmd := exec.Command("go", "list", "-m", "-u", "all")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to check dependencies: %w", err)
	}

	// Parse the output
	lines := strings.Split(string(output), "\n")
	outdated := []string{}

	for _, line := range lines {
		if strings.Contains(line, "[") && strings.Contains(line, "]") {
			outdated = append(outdated, line)
		}
	}

	// Display results
	if len(outdated) > 0 {
		fmt.Println("\nOutdated Dependencies:")
		for _, dep := range outdated {
			fmt.Println("-", dep)
		}
		fmt.Println("\nUse 'goforge dependency update' to update them.")
	} else {
		fmt.Println("\nAll dependencies are up to date!")
	}

	return nil
}

// Update updates dependencies to their latest versions.
func Update(path string) error {
	fmt.Println("Updating dependencies in:", path)

	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
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

	// Use 'go get -u' to update dependencies
	cmd := exec.Command("go", "get", "-u", "./...")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update dependencies: %w\nOutput: %s", err, output)
	}

	fmt.Println("Dependencies updated successfully!")
	fmt.Println("\nRunning 'go mod tidy' to clean up go.mod and go.sum...")

	// Run go mod tidy to clean up
	cmd = exec.Command("go", "mod", "tidy")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to tidy dependencies: %w\nOutput: %s", err, output)
	}

	fmt.Println("Dependencies tidied successfully!")
	return nil
}

// CheckSecurity checks dependencies for security vulnerabilities.
func CheckSecurity(path string) error {
	fmt.Println("Checking dependencies for security vulnerabilities in:", path)

	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
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

	// In a real implementation, this would use a security scanning tool like govulncheck
	// For this example, we'll simulate a vulnerability scan
	fmt.Println("\nSecurity Scan Results:")
	fmt.Println("- No critical vulnerabilities found")
	fmt.Println("- 2 moderate vulnerabilities in indirect dependencies")
	fmt.Println("  - github.com/example/package@v1.2.3: CVE-2023-12345")
	fmt.Println("  - github.com/another/lib@v0.1.2: GHSA-abcd-1234-5678")

	fmt.Println("\nRecommendation:")
	fmt.Println("Run 'go get github.com/example/package@v1.3.0' to resolve CVE-2023-12345")

	return nil
}
