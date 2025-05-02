#!/bin/bash
if [ $# -lt 2 ]; then
    echo "Usage: $0 <username> <password>"
    echo "Example: $0 admin securepassword"
    exit 1
fi

USERNAME=$1
PASSWORD=$2

echo "⚙️ Creating user: $USERNAME"

CID=$(docker ps --filter "name=backend" -q)
[ -z "$CID" ]

if [ -z "$CID" ]; then
    echo "❌ Backend container not found. Make sure it's running."
    exit 1
fi

echo "📦 Found backend container: $CID"

echo "📝 Preparing environment…"
grep -v '^#' /etc/environment | docker exec -i $CID tee /app/.env >/dev/null


echo "🔑 Creating user..."
docker exec -i $CID /app/dbcmd "$USERNAME" "$PASSWORD"

# Clean up the temporary .env file
docker exec $BACKEND_CONTAINER rm -f /app/.env

if [ $? -eq 0 ]; then
    echo "✅ User '$USERNAME' created successfully!"
else
    echo "❌ Failed to create user."
    exit 1
fi