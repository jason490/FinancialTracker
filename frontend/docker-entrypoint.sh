#!/bin/sh
set -e

# Keep container node_modules in sync when package.json changes on the host.
npm install

exec "$@"
