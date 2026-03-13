#!/bin/sh

# Docker entrypoint script for frontend
set -e

# Replace environment variables in built files if needed
# This allows runtime configuration of the frontend

# If API_URL is provided, replace it in the built files
if [ ! -z "$API_URL" ]; then
    echo "Configuring API URL: $API_URL"
    find /usr/share/nginx/html -name "*.js" -exec sed -i "s|__API_URL__|$API_URL|g" {} \;
fi

# If WEBSOCKET_URL is provided, replace it in the built files
if [ ! -z "$WEBSOCKET_URL" ]; then
    echo "Configuring WebSocket URL: $WEBSOCKET_URL"
    find /usr/share/nginx/html -name "*.js" -exec sed -i "s|__WEBSOCKET_URL__|$WEBSOCKET_URL|g" {} \;
fi

# Execute the main command
exec "$@"