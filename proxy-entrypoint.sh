#!/bin/sh
set -e

# Fix permissions on Caddy volumes
chown -R caddy:caddy /data /config /srv /etc/caddy

# Drop privileges and execute Caddy
exec su-exec caddy caddy run --config /etc/caddy/Caddyfile --adapter caddyfile
