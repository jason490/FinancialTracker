#!/bin/bash

# Exit on error
set -e

# Navigate to the project root (where the script is located relative to e2e/)
# This ensures we build from the correct context
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$DIR")"

echo "🚀 Building E2E Test Image..."
docker build -t financial-tracker-e2e "$DIR"

echo "🧪 Running E2E Tests..."
# Use --network host so the container shares the host's network stack.
# This is more reliable than --add-host=host.docker.internal:host-gateway on
# Linux and allows the test browser to reach the Caddy proxy on localhost:80.
docker run --rm \
    --network host \
    financial-tracker-e2e \
    pytest --base-url http://localhost "$@"

echo "✅ Tests complete."
