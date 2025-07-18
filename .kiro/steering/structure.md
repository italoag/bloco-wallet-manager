# Project Structure

## Directory Organization

```
blocowallet/
├── cmd/                    # Application entry points
│   └── blocowallet/        # Main CLI application
├── internal/               # Private application code
│   ├── blockchain/         # Blockchain interaction logic
│   ├── constants/          # Application constants
│   ├── storage/            # Database models and repositories
│   ├── ui/                 # Terminal UI components
│   └── wallet/             # Core wallet functionality
│       └── testdata/       # Test data for wallet functionality
├── keystores/              # User keystores directory
├── pkg/                    # Public packages that can be imported
│   ├── config/             # Configuration management
│   ├── localization/       # Internationalization support
│   └── logger/             # Logging utilities
└── .kiro/                  # Kiro AI assistant configuration
    └── specs/              # Feature specifications
```

## Code Organization Principles

1. **Package Structure**:
   - `cmd/`: Contains application entry points
   - `internal/`: Private application code not meant to be imported by other projects
   - `pkg/`: Public packages that can be imported by other projects

2. **Dependency Direction**:
   - Packages at higher levels can import packages from lower levels
   - Lower-level packages should not import higher-level packages
   - Example: `cmd` can import `internal` and `pkg`, but `pkg` should not import `internal`

3. **Repository Pattern**:
   - Data access is abstracted through repository interfaces
   - Implementation details are hidden behind these interfaces
   - Allows for easy switching between storage backends

4. **Service Layer**:
   - Business logic is encapsulated in service structs
   - Services depend on repositories for data access
   - UI components interact with services, not directly with repositories

5. **Error Handling**:
   - Custom error types for domain-specific errors
   - Errors include context information for better debugging
   - Localized error messages for user-facing errors

6. **Testing**:
   - Test files are placed alongside the code they test with `_test.go` suffix
   - Test data is stored in `testdata` directories
   - Mock implementations are provided for testing with `mock_` prefix