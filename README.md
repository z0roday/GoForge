# GoForge - Go Development Companion

GoForge is a comprehensive development companion for Go projects, providing a suite of productivity tools for developers in a single package. It is available as a command-line interface (CLI), API, and web interface.

## Features

- **Smart Code Analysis**: Analyze project structure, identify architectural issues, and suggest improvements
- **Dependency Management**: Automatically check and update dependencies with vulnerability detection
- **Efficient Profiling**: Visual tools for identifying performance bottlenecks
- **Docker and Container Building**: Automatic Dockerfile and Kubernetes config generation
- **Smart Testing**: Generate high-coverage tests automatically
- **Automatic Documentation**: Generate API and user documentation

## Installation

### Using Go Install

```bash
go install github.com/z0roday/goforge@latest
```

### From Source

```bash
git clone https://github.com/z0roday/goforge.git
cd goforge
go build
```

## Usage

GoForge provides a simple, intuitive command-line interface:

```bash
goforge [command] [subcommand] [options]
```

### Code Analysis

Analyze your project structure:

```bash
goforge analyze structure ./my-project
```

Analyze code quality:

```bash
goforge analyze quality ./my-project
```

### Dependency Management

Check for outdated dependencies:

```bash
goforge dependency check
```

Update dependencies:

```bash
goforge dependency update
```

Check for security vulnerabilities:

```bash
goforge dependency security
```

### Profiling

Profile CPU usage:

```bash
goforge profile cpu ./my-binary -o cpu.pprof -d 30
```

Profile memory usage:

```bash
goforge profile memory ./my-binary -o mem.pprof
```

Visualize profile data:

```bash
goforge profile visualize cpu.pprof
```

### Container Generation

Generate a Dockerfile:

```bash
goforge container dockerfile -o Dockerfile -b golang:alpine
```

Generate Kubernetes manifests:

```bash
goforge container kubernetes -o kubernetes -i myapp:latest
```

### Test Generation

Generate tests for a file or package:

```bash
goforge test generate ./pkg/mypackage -t
```

Analyze test coverage:

```bash
goforge test coverage -t 80.0 -o coverage.html
```

### Documentation Generation

Generate API documentation:

```bash
goforge docs api -o api-docs -f html
```

Generate user documentation:

```bash
goforge docs user -o user-docs -f markdown
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 