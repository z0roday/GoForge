package cmd

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

// WebCommand returns the CLI command for starting the web interface.
func WebCommand() *cli.Command {
	return &cli.Command{
		Name:  "web",
		Usage: "Start the web interface",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   "8081",
				Usage:   "Port to run the web interface on",
			},
		},
		Action: func(c *cli.Context) error {
			port := c.String("port")
			return startWebServer(port)
		},
	}
}

// startWebServer starts the web interface on the specified port.
func startWebServer(port string) error {
	fmt.Printf("Starting web interface on port %s...\n", port)

	// Create temporary directory for static files
	tempDir, err := os.MkdirTemp("", "goforge-web")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Create static files
	createStaticFiles(tempDir)

	// Define routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, filepath.Join(tempDir, "templates/index.html"), nil)
	})

	http.HandleFunc("/analyze", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, filepath.Join(tempDir, "templates/analyze.html"), nil)
	})

	http.HandleFunc("/dependency", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, filepath.Join(tempDir, "templates/dependency.html"), nil)
	})

	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, filepath.Join(tempDir, "templates/profile.html"), nil)
	})

	http.HandleFunc("/container", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, filepath.Join(tempDir, "templates/container.html"), nil)
	})

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, filepath.Join(tempDir, "templates/test.html"), nil)
	})

	http.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, filepath.Join(tempDir, "templates/docs.html"), nil)
	})

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(tempDir, "static")))))

	// Start the server
	addr := ":" + port
	fmt.Printf("Web interface is running at http://localhost%s\n", addr)
	fmt.Println("Press Ctrl+C to stop")
	return http.ListenAndServe(addr, nil)
}

// renderTemplate renders the specified template.
func renderTemplate(w http.ResponseWriter, templatePath string, data interface{}) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse template: %v", err), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute template: %v", err), http.StatusInternalServerError)
	}
}

// createStaticFiles creates the static files for the web interface.
func createStaticFiles(tempDir string) {
	// Create directories
	templatesDir := filepath.Join(tempDir, "templates")
	staticDir := filepath.Join(tempDir, "static")
	cssDir := filepath.Join(staticDir, "css")
	jsDir := filepath.Join(staticDir, "js")

	os.MkdirAll(templatesDir, 0755)
	os.MkdirAll(cssDir, 0755)
	os.MkdirAll(jsDir, 0755)

	// Create base template
	baseHTML := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GoForge - Go Development Companion</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <header>
        <div class="logo">GoForge</div>
        <nav>
            <ul>
                <li><a href="/">Home</a></li>
                <li><a href="/analyze">Analyze</a></li>
                <li><a href="/dependency">Dependencies</a></li>
                <li><a href="/profile">Profile</a></li>
                <li><a href="/container">Containers</a></li>
                <li><a href="/test">Testing</a></li>
                <li><a href="/docs">Docs</a></li>
            </ul>
        </nav>
    </header>
    <main>
        {{.Content}}
    </main>
    <footer>
        <p>&copy; 2023 GoForge - A Go Development Companion</p>
    </footer>
    <script src="/static/js/script.js"></script>
</body>
</html>
`

	// Create templates
	templates := map[string]string{
		"index.html": `
<div class="hero">
    <h1>GoForge</h1>
    <p>A comprehensive development companion for Go projects</p>
    <div class="cta-buttons">
        <a href="/analyze" class="cta-button">Analyze Code</a>
        <a href="/dependency" class="cta-button">Manage Dependencies</a>
    </div>
</div>
<div class="features">
    <div class="feature">
        <h2>Smart Code Analysis</h2>
        <p>Analyze project structure, identify architectural issues, and suggest improvements</p>
    </div>
    <div class="feature">
        <h2>Dependency Management</h2>
        <p>Automatically check and update dependencies with vulnerability detection</p>
    </div>
    <div class="feature">
        <h2>Efficient Profiling</h2>
        <p>Visual tools for identifying performance bottlenecks</p>
    </div>
    <div class="feature">
        <h2>Container Building</h2>
        <p>Automatic Dockerfile and Kubernetes config generation</p>
    </div>
    <div class="feature">
        <h2>Smart Testing</h2>
        <p>Generate high-coverage tests automatically</p>
    </div>
    <div class="feature">
        <h2>Documentation</h2>
        <p>Generate API and user documentation</p>
    </div>
</div>
`,
		"analyze.html": `
<div class="page-header">
    <h1>Code Analysis</h1>
    <p>Analyze your Go project structure and code quality</p>
</div>
<div class="tool-form">
    <form id="analyzeForm">
        <div class="form-group">
            <label for="projectPath">Project Path:</label>
            <input type="text" id="projectPath" name="path" placeholder="/path/to/your/project" required>
        </div>
        <div class="form-group">
            <label for="analysisType">Analysis Type:</label>
            <select id="analysisType" name="type">
                <option value="structure">Project Structure</option>
                <option value="quality">Code Quality</option>
            </select>
        </div>
        <button type="submit" class="submit-button">Analyze</button>
    </form>
</div>
<div id="results" class="results"></div>
`,
		"dependency.html": `
<div class="page-header">
    <h1>Dependency Management</h1>
    <p>Check and update your project dependencies</p>
</div>
<div class="tool-form">
    <form id="dependencyForm">
        <div class="form-group">
            <label for="projectPath">Project Path:</label>
            <input type="text" id="projectPath" name="path" placeholder="/path/to/your/project" required>
        </div>
        <div class="form-group">
            <label for="depAction">Action:</label>
            <select id="depAction" name="action">
                <option value="check">Check Outdated</option>
                <option value="update">Update All</option>
                <option value="security">Security Check</option>
            </select>
        </div>
        <button type="submit" class="submit-button">Run</button>
    </form>
</div>
<div id="results" class="results"></div>
`,
		"profile.html": `
<div class="page-header">
    <h1>Application Profiling</h1>
    <p>Profile your Go application performance</p>
</div>
<div class="tool-form">
    <form id="profileForm">
        <div class="form-group">
            <label for="binaryPath">Binary Path:</label>
            <input type="text" id="binaryPath" name="binary" placeholder="/path/to/your/binary" required>
        </div>
        <div class="form-group">
            <label for="profileType">Profile Type:</label>
            <select id="profileType" name="type">
                <option value="cpu">CPU Profile</option>
                <option value="memory">Memory Profile</option>
            </select>
        </div>
        <div class="form-group">
            <label for="duration">Duration (seconds):</label>
            <input type="number" id="duration" name="duration" value="30" min="5" max="300">
        </div>
        <button type="submit" class="submit-button">Profile</button>
    </form>
</div>
<div id="results" class="results"></div>
`,
		"container.html": `
<div class="page-header">
    <h1>Container Generation</h1>
    <p>Generate Docker and Kubernetes configurations</p>
</div>
<div class="tool-form">
    <form id="containerForm">
        <div class="form-group">
            <label for="projectPath">Project Path:</label>
            <input type="text" id="projectPath" name="path" placeholder="/path/to/your/project" required>
        </div>
        <div class="form-group">
            <label for="containerType">Generation Type:</label>
            <select id="containerType" name="type">
                <option value="dockerfile">Dockerfile</option>
                <option value="kubernetes">Kubernetes Manifests</option>
            </select>
        </div>
        <div class="form-group dockerfile-options">
            <label for="baseImage">Base Image:</label>
            <input type="text" id="baseImage" name="base" value="golang:alpine">
        </div>
        <div class="form-group k8s-options" style="display:none">
            <label for="imageName">Image Name:</label>
            <input type="text" id="imageName" name="image" placeholder="myapp:latest">
        </div>
        <button type="submit" class="submit-button">Generate</button>
    </form>
</div>
<div id="results" class="results"></div>
`,
		"test.html": `
<div class="page-header">
    <h1>Test Management</h1>
    <p>Generate and analyze tests for your Go project</p>
</div>
<div class="tool-form">
    <form id="testForm">
        <div class="form-group">
            <label for="projectPath">Project Path:</label>
            <input type="text" id="projectPath" name="path" placeholder="/path/to/your/project" required>
        </div>
        <div class="form-group">
            <label for="testAction">Action:</label>
            <select id="testAction" name="action">
                <option value="generate">Generate Tests</option>
                <option value="coverage">Analyze Coverage</option>
            </select>
        </div>
        <div class="form-group gen-options">
            <label for="tableTests">Table-Driven Tests:</label>
            <input type="checkbox" id="tableTests" name="table" value="true">
        </div>
        <div class="form-group coverage-options" style="display:none">
            <label for="threshold">Coverage Threshold (%):</label>
            <input type="number" id="threshold" name="threshold" value="80" min="0" max="100">
        </div>
        <button type="submit" class="submit-button">Run</button>
    </form>
</div>
<div id="results" class="results"></div>
`,
		"docs.html": `
<div class="page-header">
    <h1>Documentation Generation</h1>
    <p>Generate documentation for your Go project</p>
</div>
<div class="tool-form">
    <form id="docsForm">
        <div class="form-group">
            <label for="projectPath">Project Path:</label>
            <input type="text" id="projectPath" name="path" placeholder="/path/to/your/project" required>
        </div>
        <div class="form-group">
            <label for="docType">Documentation Type:</label>
            <select id="docType" name="type">
                <option value="api">API Documentation</option>
                <option value="user">User Documentation</option>
            </select>
        </div>
        <div class="form-group">
            <label for="docFormat">Format:</label>
            <select id="docFormat" name="format">
                <option value="html">HTML</option>
                <option value="markdown">Markdown</option>
            </select>
        </div>
        <button type="submit" class="submit-button">Generate</button>
    </form>
</div>
<div id="results" class="results"></div>
`,
	}

	// Write template files
	for name, content := range templates {
		// Inject the base template structure
		fullContent := baseHTML
		fullContent = strings.Replace(fullContent, "{{.Content}}", content, 1)

		filePath := filepath.Join(templatesDir, name)
		os.WriteFile(filePath, []byte(fullContent), 0644)
	}

	// Create CSS
	cssContent := `
body {
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    margin: 0;
    padding: 0;
    color: #333;
    line-height: 1.6;
}

header {
    background-color: #2c3e50;
    color: white;
    padding: 1rem 2rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.logo {
    font-size: 1.5rem;
    font-weight: bold;
}

nav ul {
    display: flex;
    list-style: none;
    margin: 0;
    padding: 0;
}

nav li {
    margin-left: 1.5rem;
}

nav a {
    color: white;
    text-decoration: none;
    transition: color 0.3s;
}

nav a:hover {
    color: #3498db;
}

main {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
}

footer {
    background-color: #2c3e50;
    color: white;
    text-align: center;
    padding: 1rem;
    margin-top: 2rem;
}

.hero {
    text-align: center;
    padding: 3rem 1rem;
    background-color: #f8f9fa;
    border-radius: 8px;
    margin-bottom: 2rem;
}

.hero h1 {
    font-size: 2.5rem;
    margin-bottom: 1rem;
    color: #2c3e50;
}

.cta-buttons {
    margin-top: 2rem;
}

.cta-button {
    display: inline-block;
    padding: 0.75rem 1.5rem;
    background-color: #3498db;
    color: white;
    text-decoration: none;
    border-radius: 4px;
    font-weight: bold;
    margin: 0 0.5rem;
    transition: background-color 0.3s;
}

.cta-button:hover {
    background-color: #2980b9;
}

.features {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 2rem;
}

.feature {
    background-color: white;
    border-radius: 8px;
    padding: 1.5rem;
    box-shadow: 0 2px 5px rgba(0,0,0,0.1);
    transition: transform 0.3s, box-shadow 0.3s;
}

.feature:hover {
    transform: translateY(-5px);
    box-shadow: 0 5px 15px rgba(0,0,0,0.1);
}

.feature h2 {
    color: #2c3e50;
    margin-top: 0;
}

.page-header {
    text-align: center;
    margin-bottom: 2rem;
}

.page-header h1 {
    color: #2c3e50;
}

.tool-form {
    background-color: white;
    border-radius: 8px;
    padding: 2rem;
    box-shadow: 0 2px 5px rgba(0,0,0,0.1);
    margin-bottom: 2rem;
}

.form-group {
    margin-bottom: 1.5rem;
}

.form-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: bold;
}

.form-group input[type="text"],
.form-group input[type="number"],
.form-group select {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
}

.submit-button {
    background-color: #3498db;
    color: white;
    border: none;
    padding: 0.75rem 1.5rem;
    font-size: 1rem;
    font-weight: bold;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.3s;
}

.submit-button:hover {
    background-color: #2980b9;
}

.results {
    background-color: #f8f9fa;
    border-radius: 8px;
    padding: 1.5rem;
    white-space: pre-wrap;
    font-family: monospace;
}
`
	os.WriteFile(filepath.Join(cssDir, "style.css"), []byte(cssContent), 0644)

	// Create JavaScript
	jsContent := `
document.addEventListener('DOMContentLoaded', function() {
    // Form submission handlers
    const forms = {
        'analyzeForm': '/api/analyze/',
        'dependencyForm': '/api/dependency/',
        'profileForm': '/api/profile/',
        'containerForm': '/api/container/',
        'testForm': '/api/test/',
        'docsForm': '/api/docs/generate'
    };

    for (const [formId, apiEndpoint] of Object.entries(forms)) {
        const form = document.getElementById(formId);
        if (form) {
            form.addEventListener('submit', function(e) {
                e.preventDefault();
                const formData = new FormData(form);
                const resultsDiv = document.getElementById('results');
                
                resultsDiv.textContent = 'Processing request...';
                
                // Custom handling based on form type
                let endpoint = apiEndpoint;
                if (formId === 'analyzeForm') {
                    const type = formData.get('type');
                    endpoint += type;
                }
                
                fetch(endpoint, {
                    method: 'POST',
                    body: formData
                })
                .then(response => response.json())
                .then(data => {
                    if (data.error) {
                        resultsDiv.textContent = 'Error: ' + data.error;
                    } else {
                        if (data.data && data.data.output) {
                            resultsDiv.textContent = data.data.output;
                        } else {
                            resultsDiv.textContent = JSON.stringify(data, null, 2);
                        }
                    }
                })
                .catch(error => {
                    resultsDiv.textContent = 'Error: ' + error.message;
                });
            });
        }
    }

    // Dynamic form controls
    const setupDynamicFormControls = () => {
        // Container form
        const containerType = document.getElementById('containerType');
        if (containerType) {
            containerType.addEventListener('change', function() {
                const dockerfileOptions = document.querySelector('.dockerfile-options');
                const k8sOptions = document.querySelector('.k8s-options');
                
                if (this.value === 'dockerfile') {
                    dockerfileOptions.style.display = 'block';
                    k8sOptions.style.display = 'none';
                } else {
                    dockerfileOptions.style.display = 'none';
                    k8sOptions.style.display = 'block';
                }
            });
        }

        // Test form
        const testAction = document.getElementById('testAction');
        if (testAction) {
            testAction.addEventListener('change', function() {
                const genOptions = document.querySelector('.gen-options');
                const coverageOptions = document.querySelector('.coverage-options');
                
                if (this.value === 'generate') {
                    genOptions.style.display = 'block';
                    coverageOptions.style.display = 'none';
                } else {
                    genOptions.style.display = 'none';
                    coverageOptions.style.display = 'block';
                }
            });
        }
    };

    setupDynamicFormControls();
});
`
	os.WriteFile(filepath.Join(jsDir, "script.js"), []byte(jsContent), 0644)
}
