#!/bin/sh
set -e

# Fix permissions on the mounted SQLite data volume
# This ensures existing databases created by root are accessible
chown -R appuser:appgroup /data

# Drop privileges and execute the Go application
exec su-exec appuser /app/main
