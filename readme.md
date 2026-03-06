# CRONEX

A simple, lightweight cron job scheduler written in Go that executes scheduled tasks based on JSON configuration.

## Features

- 🕐 **Cron scheduling with second precision** - Uses extended cron format with seconds support
- 📝 **JSON-based configuration** - Easy to configure and manage scheduled tasks
- 🚀 **Cross-platform** - Works on Windows and Linux
- 💬 **Human-readable schedules** - Converts cron expressions to readable descriptions
- 📊 **Real-time output** - Displays task execution output immediately
- ⚙️ **Enable/disable tasks** - Control which tasks run without removing configuration

## Installation

### Download Pre-built Binaries

Download the latest binaries from the [GitHub Releases](../../releases) page:
- `cronex-windows-x64.exe` - Windows x64
- `cronex-linux-x64` - Linux x64

### Build from Source

Requirements:
- Go 1.25 or later

```bash
git clone https://github.com/dunglex/cronex.git
cd cronex
go build -o cronex .
```

## Configuration

Create a `cron.json` file with your scheduled tasks:

```json
[
  {
    "name": "task-name",
    "cron": "*/10 * * * * *",
    "cmd": "command",
    "args": ["arg1", "arg2"],
    "enabled": true
  }
]
```

### Configuration Fields

- **name**: Descriptive name for the task
- **cron**: Cron expression with seconds (6 fields: `second minute hour day month weekday`)
- **cmd**: Command to execute
- **args**: Array of command arguments (optional)
- **enabled**: Whether the task is active (true/false)

### Cron Expression Format

The cron format includes seconds as the first field:

```
* * * * * *
│ │ │ │ │ │
│ │ │ │ │ └─── Day of week (0-6, Sunday=0)
│ │ │ │ └───── Month (1-12)
│ │ │ └─────── Day of month (1-31)
│ │ └───────── Hour (0-23)
│ └─────────── Minute (0-59)
└───────────── Second (0-59)
```

### Example Schedules

- `*/10 * * * * *` - Every 10 seconds
- `0 */5 * * * *` - Every 5 minutes
- `0 0 */2 * * *` - Every 2 hours
- `0 30 9 * * *` - Daily at 9:30 AM
- `0 0 0 * * 0` - Every Sunday at midnight

## Usage

Run with the default configuration file (`./cron.json`):

```bash
./cronex
```

Specify a custom configuration file:

```bash
./cronex -config /path/to/config.json
```

## Example Configuration

```json
[
  {
    "name": "memory-check",
    "cron": "*/10 * * * * *",
    "cmd": "powershell",
    "args": ["-NoProfile", "-Command", "Get-Process | Measure-Object -Property WorkingSet -Sum"],
    "enabled": true
  },
  {
    "name": "backup",
    "cron": "0 0 2 * * *",
    "cmd": "bash",
    "args": ["-c", "tar -czf backup.tar.gz /data"],
    "enabled": true
  }
]
```

## Building Releases

Releases are automatically built for Windows x64 and Linux x64 when a new tag is pushed:

```bash
git tag v1.0.0
git push origin v1.0.0
```

The GitHub Actions workflow will build both binaries and make them available as artifacts.

## License

See [LICENSE](LICENSE) file for details.