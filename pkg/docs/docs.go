package docs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

// UserDocTemplate is a template for generating basic user documentation.
const UserDocTemplate = `# User Guide for {{.AppName}}

## Introduction

This document provides information on how to use the {{.AppName}} application effectively.

## Installation

To install {{.AppName}}, run:

` + "```" + `
go install github.com/yourusername/{{.AppName}}@latest
` + "```" + `

## Usage

{{.AppName}} provides the following commands:

` + "```" + `
{{.AppName}} [command] [options]
` + "```" + `

### Available Commands

- **analyze**: Analyze Go code structure and quality
- **dependency**: Manage project dependencies
- **profile**: Profile application performance
- **container**: Generate container configurations
- **test**: Test management utilities
- **docs**: Generate documentation

## Examples

### Analyzing Code

To analyze your project structure:

` + "```" + `
{{.AppName}} analyze structure ./my-project
` + "```" + `

To analyze code quality:

` + "```" + `
{{.AppName}} analyze quality ./my-project
` + "```" + `

### Managing Dependencies

To check for outdated dependencies:

` + "```" + `
{{.AppName}} dependency check
` + "```" + `

To update dependencies:

` + "```" + `
{{.AppName}} dependency update
` + "```" + `

## Support

For support, please open an issue on the GitHub repository.
`

// UserDocData holds data for the user documentation template.
type UserDocData struct {
	AppName string
}

// GenerateAPIDoc generates API documentation for a Go project.
func GenerateAPIDoc(path string, outputDir string, format string) error {
	fmt.Printf("Generating API documentation for %s in %s format\n", path, format)

	// Get absolute paths
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	absOutput, err := filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for output: %w", err)
	}

	// Create output directory if it doesn't exist
	err = os.MkdirAll(absOutput, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// For HTML format, use go doc -html
	if format == "html" {
		// Save current directory
		originalDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		defer os.Chdir(originalDir)

		// Change to project directory
		err = os.Chdir(absPath)
		if err != nil {
			return fmt.Errorf("failed to change to project directory: %w", err)
		}

		// Create index.html
		indexPath := filepath.Join(absOutput, "index.html")
		cmd := exec.Command("go", "doc", "-html", "./...")
		indexFile, err := os.Create(indexPath)
		if err != nil {
			return fmt.Errorf("failed to create index.html: %w", err)
		}
		defer indexFile.Close()

		cmd.Stdout = indexFile
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to generate HTML documentation: %w", err)
		}

		fmt.Printf("API documentation generated at: %s\n", indexPath)
	} else if format == "markdown" {
		// For markdown format, use go doc
		packages, err := filepath.Glob(filepath.Join(absPath, "pkg", "*"))
		if err != nil {
			return fmt.Errorf("failed to list packages: %w", err)
		}

		// Create index file
		indexPath := filepath.Join(absOutput, "README.md")
		indexFile, err := os.Create(indexPath)
		if err != nil {
			return fmt.Errorf("failed to create README.md: %w", err)
		}
		defer indexFile.Close()

		fmt.Fprintln(indexFile, "# API Documentation\n")
		fmt.Fprintln(indexFile, "## Packages\n")

		// Document each package
		for _, pkg := range packages {
			pkgName := filepath.Base(pkg)
			fmt.Fprintf(indexFile, "- [%s](%s.md)\n", pkgName, pkgName)

			// Generate documentation for the package
			pkgDocPath := filepath.Join(absOutput, pkgName+".md")
			pkgDocFile, err := os.Create(pkgDocPath)
			if err != nil {
				return fmt.Errorf("failed to create package documentation file: %w", err)
			}

			pkgImportPath := fmt.Sprintf("./pkg/%s", pkgName)
			cmd := exec.Command("go", "doc", "-all", pkgImportPath)
			cmd.Stdout = pkgDocFile
			err = cmd.Run()
			pkgDocFile.Close()
			if err != nil {
				return fmt.Errorf("failed to generate documentation for package %s: %w", pkgName, err)
			}
		}

		fmt.Printf("API documentation generated at: %s\n", absOutput)
	} else {
		return fmt.Errorf("unsupported format: %s (supported: html, markdown)", format)
	}

	return nil
}

// GenerateUserDoc generates user documentation for a Go project.
func GenerateUserDoc(path string, outputDir string, format string) error {
	fmt.Printf("Generating user documentation for %s in %s format\n", path, format)

	// Get absolute paths
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	absOutput, err := filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for output: %w", err)
	}

	// Create output directory if it doesn't exist
	err = os.MkdirAll(absOutput, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Determine app name from directory
	appName := filepath.Base(absPath)

	// Create template data
	data := UserDocData{
		AppName: appName,
	}

	// Parse and execute the template
	tmpl, err := template.New("userdoc").Parse(UserDocTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse user doc template: %w", err)
	}

	// Create markdown file
	mdPath := filepath.Join(absOutput, "user-guide.md")
	mdFile, err := os.Create(mdPath)
	if err != nil {
		return fmt.Errorf("failed to create user guide file: %w", err)
	}
	defer mdFile.Close()

	// Execute the template
	err = tmpl.Execute(mdFile, data)
	if err != nil {
		return fmt.Errorf("failed to execute user doc template: %w", err)
	}

	fmt.Printf("User documentation markdown generated at: %s\n", mdPath)

	// If HTML format is requested, convert markdown to HTML
	if format == "html" {
		// Check if pandoc is available (simplistic check)
		_, err := exec.LookPath("pandoc")
		if err != nil {
			fmt.Println("WARNING: pandoc not found, cannot convert to HTML. Using markdown instead.")
			return nil
		}

		// Convert markdown to HTML using pandoc
		htmlPath := filepath.Join(absOutput, "user-guide.html")
		cmd := exec.Command("pandoc", "-s", mdPath, "-o", htmlPath)
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to convert markdown to HTML: %w", err)
		}

		fmt.Printf("User documentation HTML generated at: %s\n", htmlPath)
	}

	return nil
}
