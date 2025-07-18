# Configuration Management System

This package provides a robust configuration management system for the BLOCO Wallet Manager, focusing on safe file operations, atomic writes, and proper error handling.

## Components

### ConfigFileManager

The `ConfigFileManager` handles all file operations related to configuration files, ensuring data integrity and providing recovery mechanisms.

```go
type ConfigFileManager struct {
    configPath string
}
```

#### Key Features

- **Atomic Write Operations**: Uses a write-to-temporary-file and rename approach to prevent configuration corruption
- **Automatic Backup System**: Creates timestamped backups before any configuration changes
- **Restore Capability**: Can restore from backups in case of write failures
- **Directory Management**: Ensures configuration directories exist before operations
- **Validation**: Provides basic validation of configuration files

#### Usage Example

```go
// Create a new config file manager
cfgManager := config.NewConfigFileManager("/path/to/config.toml")

// Read configuration
lines, err := cfgManager.ReadConfig()
if err != nil {
    // Handle error
}

// Modify configuration
// ...

// Write configuration back safely
err = cfgManager.WriteConfig(lines)
if err != nil {
    // Handle error
}
```

### Error Handling

The configuration management system provides detailed error messages with proper error wrapping to help diagnose issues:

- File not found errors
- Read/write permission errors
- Directory creation errors
- Backup and restore errors

### Backup System

Backups are automatically created before any write operation:

- Backups are stored with a timestamp in the filename: `config.toml.20250718_120000.bak`
- In case of write failure, the system attempts to restore from the backup
- Backups can be manually restored using the `RestoreConfig` method

## Integration with Network Configuration

This configuration management system is part of a larger effort to fix network configuration formatting issues in the BLOCO Wallet Manager. It works alongside:

1. **TOML Section Manager**: Handles finding, removing, and adding sections in TOML files
2. **Network Key Generation**: Creates valid, unique keys for network configurations
3. **UI Integration**: Provides proper feedback to users during configuration operations

## Implementation Details

### File Operations

- **Reading**: Files are read entirely into memory and split into lines
- **Writing**: Content is joined with newlines and written atomically
- **Validation**: Basic validation checks if the file exists and can be read

### Security Considerations

- File permissions are set to 0644 (read/write for owner, read for others)
- Temporary files use the same name with a `.tmp` suffix
- Backups are stored in the same directory as the original file

## Future Improvements

- Enhanced TOML validation using a proper TOML parser
- Support for different configuration formats (JSON, YAML)
- Encryption options for sensitive configuration data
- Compression for large configuration files or backups