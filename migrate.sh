#!/bin/bash
DB_CONTAINER=$(docker ps --filter "name=db" --format "{{.Names}}")

if [ -z "$DB_CONTAINER" ]; then
    echo "Error: PostgreSQL container not found"
    exit 1
fi

echo "Found database container: $DB_CONTAINER"
echo "Running migrations from schema.sql..."

docker exec -i $DB_CONTAINER psql -U $POSTGRES_USER -d $POSTGRES_DB -f - < ~/app/migrations/schema.sql

if [ $? -eq 0 ]; then
    echo "✅ Schema migration completed successfully"
else
    echo "❌ Schema migration failed"
    exit 1
fi

# Seed is not needed
# if [ -f migrations/seed.sql ]; then
#   echo "Running seed data..."
#   docker exec -i $DB_CONTAINER psql -U $POSTGRES_USER -d $POSTGRES_DB -f - < migrations/seed.sql

#   if [ $? -eq 0 ]; then
#     echo "✅ Seed data loaded successfully"
#   else
#     echo "❌ Seed data loading failed"
#     exit 1
#   fi
# fi

echo "Database migration complete!"
