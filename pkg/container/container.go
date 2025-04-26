package container

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// DockerfileTemplate is a template for generating a basic Dockerfile for Go applications.
const DockerfileTemplate = `FROM {{ .BaseImage }} as builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Use a small image for the final stage
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/app .

# Expose port if needed
EXPOSE 8080

# Command to run
CMD ["./app"]
`

// K8sDeploymentTemplate is a template for generating a basic Kubernetes deployment.
const K8sDeploymentTemplate = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .AppName }}
  labels:
    app: {{ .AppName }}
spec:
  replicas: 3
  selector:
    matchLabels:
      app: {{ .AppName }}
  template:
    metadata:
      labels:
        app: {{ .AppName }}
    spec:
      containers:
      - name: {{ .AppName }}
        image: {{ .Image }}
        ports:
        - containerPort: 8080
        resources:
          limits:
            cpu: "500m"
            memory: "512Mi"
          requests:
            cpu: "100m"
            memory: "128Mi"
`

// K8sServiceTemplate is a template for generating a basic Kubernetes service.
const K8sServiceTemplate = `apiVersion: v1
kind: Service
metadata:
  name: {{ .AppName }}
spec:
  selector:
    app: {{ .AppName }}
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
`

// DockerfileData holds data for the Dockerfile template.
type DockerfileData struct {
	BaseImage string
}

// K8sData holds data for the Kubernetes templates.
type K8sData struct {
	AppName string
	Image   string
}

// GenerateDockerfile creates a Dockerfile for a Go application.
func GenerateDockerfile(path string, outputFile string, baseImage string) error {
	fmt.Println("Generating Dockerfile for project at:", path)

	// Get absolute paths
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	absOutput, err := filepath.Abs(outputFile)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for output: %w", err)
	}

	// Determine app name from directory
	appName := filepath.Base(absPath)

	// Create template data
	data := DockerfileData{
		BaseImage: baseImage,
	}

	// Parse and execute the template
	tmpl, err := template.New("dockerfile").Parse(DockerfileTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse Dockerfile template: %w", err)
	}

	// Create output file
	file, err := os.Create(absOutput)
	if err != nil {
		return fmt.Errorf("failed to create Dockerfile: %w", err)
	}
	defer file.Close()

	// Execute the template
	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("failed to execute Dockerfile template: %w", err)
	}

	fmt.Printf("Dockerfile generated at: %s\n", absOutput)
	fmt.Println("\nTo build the Docker image, run:")
	fmt.Printf("docker build -t %s:latest -f %s %s\n", strings.ToLower(appName), outputFile, path)

	return nil
}

// GenerateKubernetesManifests creates Kubernetes manifests for a Go application.
func GenerateKubernetesManifests(path string, outputDir string, image string) error {
	fmt.Println("Generating Kubernetes manifests for project at:", path)

	// Get absolute paths
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	absOutput, err := filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for output: %w", err)
	}

	// Determine app name from directory
	appName := filepath.Base(absPath)

	// Use app name as image if not specified
	if image == "" {
		image = strings.ToLower(appName) + ":latest"
	}

	// Create template data
	data := K8sData{
		AppName: appName,
		Image:   image,
	}

	// Create output directory if it doesn't exist
	err = os.MkdirAll(absOutput, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate deployment manifest
	deploymentPath := filepath.Join(absOutput, "deployment.yaml")
	deploymentFile, err := os.Create(deploymentPath)
	if err != nil {
		return fmt.Errorf("failed to create deployment manifest: %w", err)
	}
	defer deploymentFile.Close()

	deploymentTmpl, err := template.New("deployment").Parse(K8sDeploymentTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse deployment template: %w", err)
	}

	err = deploymentTmpl.Execute(deploymentFile, data)
	if err != nil {
		return fmt.Errorf("failed to execute deployment template: %w", err)
	}

	// Generate service manifest
	servicePath := filepath.Join(absOutput, "service.yaml")
	serviceFile, err := os.Create(servicePath)
	if err != nil {
		return fmt.Errorf("failed to create service manifest: %w", err)
	}
	defer serviceFile.Close()

	serviceTmpl, err := template.New("service").Parse(K8sServiceTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse service template: %w", err)
	}

	err = serviceTmpl.Execute(serviceFile, data)
	if err != nil {
		return fmt.Errorf("failed to execute service template: %w", err)
	}

	fmt.Printf("Kubernetes manifests generated in: %s\n", absOutput)
	fmt.Println("\nTo apply the manifests, run:")
	fmt.Printf("kubectl apply -f %s\n", absOutput)

	return nil
}
