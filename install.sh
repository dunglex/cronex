#!/bin/bash
set -e

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root (use sudo)"
   exit 1
fi

# Create cronex user and group
echo "Creating cronex user..."
if ! id -u cronex > /dev/null 2>&1; then
    useradd --system --no-create-home --shell /bin/false cronex
fi

# Create configuration directory
echo "Setting up /etc/cronex directory..."
mkdir -p /etc/cronex
chown cronex:cronex /etc/cronex
chmod 755 /etc/cronex

# Copy binary to /usr/local/bin
echo "Installing cronex binary..."
cp cronex /usr/local/bin/cronex
chmod 755 /usr/local/bin/cronex

# Copy service file
echo "Installing systemd service..."
cp cronex.service /etc/systemd/system/cronex.service
chmod 644 /etc/systemd/system/cronex.service

# Copy config file if it exists
if [ -f "cron.json" ]; then
    echo "Copying cron.json configuration..."
    cp cron.json /etc/cronex/cron.json
    chown cronex:cronex /etc/cronex/cron.json
    chmod 640 /etc/cronex/cron.json
else
    echo "Warning: cron.json not found. Please copy your config to /etc/cronex/cron.json"
fi

# Reload systemd daemon
echo "Reloading systemd daemon..."
systemctl daemon-reload

# Enable the service
echo "Enabling cronex service..."
systemctl enable cronex

echo ""
echo "✓ Installation complete!"
echo ""
echo "Next steps:"
echo "1. Edit your configuration: sudo nano /etc/cronex/cron.json"
echo "2. Start the service: sudo systemctl start cronex"
echo "3. Check status: sudo systemctl status cronex"
echo "4. View logs: sudo journalctl -u cronex -f"
