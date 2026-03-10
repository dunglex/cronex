#!/bin/bash
set -e

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root (use sudo)"
   exit 1
fi

echo "Uninstalling cronex..."

# Stop the service
echo "Stopping cronex service..."
systemctl stop cronex || true
systemctl disable cronex || true

# Remove service file
echo "Removing systemd service file..."
rm -f /etc/systemd/system/cronex.service

# Remove binary
echo "Removing cronex binary..."
rm -f /usr/local/bin/cronex

# Remove configuration directory (optional - commented out for safety)
# echo "Removing /etc/cronex directory..."
# rm -rf /etc/cronex

# Remove cronex user
echo "Removing cronex user..."
userdel cronex || true

# Reload systemd daemon
echo "Reloading systemd daemon..."
systemctl daemon-reload

echo ""
echo "✓ Uninstall complete!"
echo ""
echo "Note: Configuration directory /etc/cronex was not removed."
echo "To remove it manually: sudo rm -rf /etc/cronex"
