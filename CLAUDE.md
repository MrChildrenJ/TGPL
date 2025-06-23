# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a learning repository containing Go programming exercises and examples from "The Go Programming Language" book. The codebase is organized by chapters with example programs demonstrating various Go concepts and features.

## Project Structure

- **Ch01/**: Basic Go programs including command-line arguments, HTTP servers, and file operations
- **Ch02/**: Type declarations, packages, and conversion utilities (includes tempconv package)
- **Ch03/**: Advanced graphics and web servers with SVG generation
- **go.mod**: Module definition with dependency on gopl.io

## Common Commands

### Running Go Programs
```bash
# Run a specific Go file
go run Ch01/go1_1.go

# Run with arguments
go run Ch01/go1_1.go arg1 arg2

# Run HTTP server examples
go run Ch03/go3_4.go  # Starts server on localhost:8000
```

### Building and Development
```bash
# Build all packages
go build ./...

# Format code
go fmt ./...

# Vet code for common mistakes
go vet ./...

# Test packages
go test ./...

# Run a single test file
go test -run TestName ./path/to/package
```

### Module Management
```bash
# Download dependencies
go mod download

# Tidy up dependencies
go mod tidy

# Verify dependencies
go mod verify
```

## Architecture Notes

### Package Organization
- Individual Go files contain standalone programs with `main()` functions
- Some have renamed main functions (e.g., `main1()`) to avoid conflicts when multiple files are in the same directory
- The `tempconv` package in Ch02 demonstrates proper Go package structure with custom types and methods

### Web Server Examples
- Ch03 contains HTTP server implementations that generate SVG graphics
- Surface plotting functionality with 3D mathematical visualization
- Color mapping based on height values for visual representation

### Key Patterns
- Most programs are educational examples focusing on specific Go language features
- Files often contain mathematical computations and visual output generation
- Error handling patterns include checking for NaN and Inf values in mathematical calculations