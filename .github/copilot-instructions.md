# Copilot Instructions for gookit/config

## Repository Overview

**gookit/config** is a comprehensive Go configuration management library that supports multiple formats (JSON, YAML, TOML, INI, HCL, ENV, Flags). It provides features like environment variable parsing, data merging, struct binding, event hooks, and remote configuration loading.

- **Language**: Go (requires Go 1.19+)
- **Size**: ~40 Go source files across multiple format driver packages
- **Type**: Library package (not an executable application)
- **Module Path**: `github.com/gookit/config/v2`
- **Coverage**: >95% test coverage with comprehensive unit tests

## Build and Test Instructions

### Prerequisites
- Go 1.19 or higher (CI tests against 1.19, 1.20, 1.21, 1.22, 1.23, 1.24)
- No additional build tools required (pure Go project)

### Dependencies Management
**ALWAYS run this first before any other operations:**
```bash
go mod tidy
```

### Building
Build the library to verify compilation:
```bash
go build ./...
```
**Expected**: No output indicates successful build. This is a library, so no binaries are generated.

### Testing
Run all tests with coverage (recommended):
```bash
go test -cover ./...
```
**Expected**: All tests pass with >95% coverage. Typical runtime: ~3 seconds.

Run tests for specific package:
```bash
go test -cover .                    # Main package only
go test -cover ./yaml               # YAML driver only  
go test -cover ./json5              # JSON5 driver only
```

**Note**: Some driver packages may show "no statements" for coverage as they only import/register drivers.

### Linting
Install and run linting tools used in CI:

```bash
# Install linting tools (one-time setup)
go install honnef.co/go/tools/cmd/staticcheck@latest
go install github.com/mgechev/revive@latest

# Run static analysis (fails CI if errors found)
staticcheck ./...

# Run style linting (excludes examples and testdata)
revive -exclude ./_examples/... -exclude ./testdata/... ./...
```

**Expected Issues**: The linting tools may report some issues that are acceptable in this codebase (unused parameters in incomplete HCL implementations, etc.). Only new issues in your changes need to be addressed.

### Running Examples
Test functionality with examples:
```bash
# Run YAML example (demonstrates typical usage)
go run _examples/yaml.go

# All examples are in _examples/ directory
ls _examples/
```

## Project Architecture and Layout

### Core Files (Repository Root)
- `config.go` - Main Config struct and core functionality
- `driver.go` - Driver interface and JSON driver implementation  
- `options.go` - Configuration options and hooks
- `load.go` - File/data loading functionality
- `read.go` - Data reading and type conversion methods
- `write.go` - Data writing and modification methods
- `export.go` - Data export and struct binding
- `util.go` - Utility functions

### Format Driver Packages
Each format has its own package with minimal code (usually just driver registration):
- `json/` - Enhanced JSON driver with additional features
- `json5/` - JSON5 format support
- `yaml/` - YAML v2 support  
- `yamlv3/` - YAML v3 support (recommended)
- `toml/` - TOML format support
- `ini/` - INI format support
- `hcl/` - HCL v1 support (basic)
- `hclv2/` - HCL v2 support (incomplete/experimental)
- `properties/` - Java properties format
- `other/` - Empty driver package for examples

### Test and Example Structure
- `*_test.go` - Unit tests alongside source files
- `issues_test.go` - Regression tests for reported issues
- `testdata/` - Test configuration files in various formats
- `_examples/` - Usage examples for each format

### Configuration Files
- `go.mod`/`go.sum` - Go module dependencies
- `.github/workflows/` - CI/CD workflows
  - `go.yml` - Unit tests on multiple Go versions
  - `lint.yml` - Code linting with revive and staticcheck  
  - `codeql.yml` - Security code scanning
  - `release.yml` - Automated releases

## CI/CD and Validation

### GitHub Actions Workflows
The repository runs comprehensive validation on every PR and push:

1. **Unit Tests** (`go.yml`): Tests on Go 1.19-1.24, generates coverage reports
2. **Linting** (`lint.yml`): Runs revive and staticcheck for code quality
3. **CodeQL** (`codeql.yml`): Security vulnerability scanning
4. **Release** (`release.yml`): Automated versioning and releases

### Pre-commit Validation
Before committing, run the same checks as CI:
```bash
# Run all tests
go test -cover ./...

# Run linting  
staticcheck ./...
revive -exclude ./_examples/... -exclude ./testdata/... ./...

# Verify build
go build ./...
```

### Test Data Requirements
- Always test new features with files in `testdata/` directory
- Each format should have `format_base.ext` and `format_other.ext` test files
- Examples in `_examples/` must work from repository root: `go run _examples/name.go`

## Key Architecture Patterns

### Driver System
The library uses a driver pattern for format support:
- `Driver` interface defines encode/decode for formats
- `StdDriver` provides standard implementation
- Each format package registers its driver on import
- JSON driver is built-in, others are optional

### Configuration Loading
Standard usage pattern:
```go
config.AddDriver(yamlv3.Driver)  // Add format support
config.LoadFiles("config.yml")   // Load config files
config.String("key")             // Read values
```

### Options and Hooks
- Options control parsing behavior (env vars, defaults, caching)
- Hook functions fire on data changes (useful for file watching)
- Set options before loading any data

### Error Handling
- Methods return errors for validation
- Panic on programmer errors (e.g., setting options after loading data)
- Default values provided for missing configuration

## Common Development Patterns

### Adding New Format Support
1. Create new driver package (copy existing as template)
2. Implement `Driver` interface with encode/decode
3. Add test files to `testdata/`
4. Create example in `_examples/`
5. Add import to main package for registration

### Adding Configuration Options
1. Add field to `Options` struct in `options.go`
2. Create option function (e.g., `func NewOption(opts *Options) { opts.Field = value }`)
3. Update `newDefaultOption()` with sensible default
4. Add tests in `config_test.go`

### Modifying Core Functionality  
- Main logic in `config.go`, `load.go`, `read.go`, `write.go`, `export.go`
- Always maintain backward compatibility
- Add comprehensive tests for new features
- Update examples if behavior changes

## Trust These Instructions

These instructions have been validated by:
- Running all build and test commands successfully  
- Executing examples to verify functionality
- Testing linting tools and CI workflows
- Examining actual CI configurations and workflows

Only search for additional information if:
- Instructions are incomplete for your specific task
- Commands fail with unexpected errors  
- Working with experimental features (like HCL v2)
- Need details about specific internal implementation

**Always run the validation commands before submitting changes to avoid CI failures.**