# Release Log

## [1.0.1] - 2026-03-10

### Added
- YAML configuration file support (.yaml and .yml extensions)
- Automatic format detection based on file extension (JSON, YAML, or YML)
- Comprehensive YAML test suite (yaml_test.go)
- YAML vs JSON parity testing
- Example YAML configuration files (cron.example.yaml, cron.test.yaml)
- Enhanced test coverage for configuration parsing


## [1.0.0] - 2026-03-06

### Added
- Initial release of CRONEX
- Cron job scheduling with second-precision using extended 6-field cron format
- JSON-based configuration file support
- Cross-platform support for Windows and Linux
- Command execution with argument support
- Enable/disable toggle for individual tasks
- Real-time command output display
- Human-readable cron expression descriptions
- Custom configuration file path via `-config` flag
- GitHub Actions workflow for automated builds (Windows x64 and Linux x64)