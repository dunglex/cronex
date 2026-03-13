# CRONEX

A simple, lightweight cron job scheduler written in Go that executes scheduled tasks based on JSON configuration.

## Features

- 🕐 **Cron scheduling with second precision** - Uses extended cron format with seconds support
- 📝 **JSON and YAML configuration** - Easy to configure and manage scheduled tasks in your preferred format
- 🚀 **Cross-platform** - Works on Windows and Linux
- 💬 **Human-readable schedules** - Converts cron expressions to readable descriptions
- 📊 **Real-time output** - Displays task execution output immediately
- ⚙️ **Enable/disable tasks** - Control which tasks run without removing configuration

## Installation

### Download Pre-built Binaries

Download the latest binaries from the [GitHub Releases](../../releases) page:
- `cronex-windows-x64.exe` - Windows x64
- `cronex-linux-x64` - Linux x64
- `cronex-linux-arm64` - Linux ARM64

### Build from Source

Requirements:
- Go 1.25 or later

```bash
git clone https://github.com/dunglex/cronex.git
cd cronex
go build -o cronex .
```

### Install as Linux Systemd Service

For Linux users, you can install cronex as a systemd service to run automatically on startup:

```bash
sudo ./install.sh
```

This installation script will:
- Build the binary
- Create a dedicated `cronex` user and group
- Install the binary to `/usr/local/bin/`
- Install the systemd service file
- Set up the configuration directory at `/etc/cronex/`
- Copy your `cron.json` to the systemd configuration location

After installation:

```bash
# Edit your configuration
sudo nano /etc/cronex/cron.json

# Start the service
sudo systemctl start cronex

# Check status
sudo systemctl status cronex

# View logs
sudo journalctl -u cronex -f

# Enable to run on boot (already done by install script)
sudo systemctl enable cronex
```

### Manual Systemd Setup

If you prefer to set up manually:

1. Create a system user:
   ```bash
   sudo useradd --system --no-create-home --shell /bin/false cronex
   ```

2. Copy the binary:
   ```bash
   sudo cp cronex /usr/local/bin/
   ```

3. Create config directory:
   ```bash
   sudo mkdir -p /etc/cronex
   sudo cp cron.json /etc/cronex/
   sudo chown cronex:cronex /etc/cronex
   ```

4. Install the service file:
   ```bash
   sudo cp cronex.service /etc/systemd/system/
   sudo systemctl daemon-reload
   sudo systemctl enable cronex
   sudo systemctl start cronex
   ```

## Configuration

CRONEX supports both JSON and YAML configuration formats. The format is automatically detected based on the file extension (`.json`, `.yaml`, or `.yml`).

### JSON Configuration

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

### YAML Configuration

Or create a `cron.yaml` (or `cron.yml`) file:

```yaml
- name: task-name
  cron: "*/10 * * * * *"
  cmd: command
  args:
    - arg1
    - arg2
  enabled: true
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

Specify a custom configuration file (supports both JSON and YAML):

```bash
./cronex -config /path/to/config.json
./cronex -config /path/to/config.yaml
```

## Example Configuration

### JSON Example

### JSON Format

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

### YAML Format

```yaml
- name: memory-check
  cron: "*/10 * * * * *"
  cmd: powershell
  args:
    - -NoProfile
    - -Command
    - Get-Process | Measure-Object -Property WorkingSet -Sum
  enabled: true

- name: backup
  cron: "0 0 2 * * *"
  cmd: bash
  args:
    - -c
    - tar -czf backup.tar.gz /data
  enabled: true
```

### YAML Example

```yaml
- name: memory-check
  cron: "*/10 * * * * *"
  cmd: powershell
  args:
    - "-NoProfile"
    - "-Command"
    - "Get-Process | Measure-Object -Property WorkingSet -Sum"
  enabled: true

- name: backup
  cron: "0 0 2 * * *"
  cmd: bash
  args:
    - "-c"
    - "tar -czf backup.tar.gz /data"
  enabled: true
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