#!/bin/bash
if [ $# -lt 2 ]; then
    echo "Usage: $0 <username> <password>"
    echo "Example: $0 admin securepassword"
    exit 1
fi

USERNAME=$1
PASSWORD=$2

echo "⚙️ Creating user: $USERNAME"

BACKEND_CONTAINER=$(docker ps --filter "name=backend" --format "{{.ID}}")

if [ -z "$BACKEND_CONTAINER" ]; then
    echo "❌ Backend container not found. Make sure it's running."
    exit 1
fi

echo "📦 Found backend container: $BACKEND_CONTAINER"

echo "📝 Preparing environment..."
docker exec $BACKEND_CONTAINER bash -c "grep -v '^#' /etc/environment > /app/.env"

echo "🔑 Creating user..."
docker exec $BACKEND_CONTAINER go run cmd/db/main.go "$USERNAME" "$PASSWORD"

# Clean up the temporary .env file
docker exec $BACKEND_CONTAINER rm -f /app/.env

if [ $? -eq 0 ]; then
    echo "✅ User '$USERNAME' created successfully!"
else
    echo "❌ Failed to create user."
    exit 1
fi