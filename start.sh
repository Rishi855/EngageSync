#!/bin/sh

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to start..."
until PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c '\q'; do
  echo "Postgres is unavailable - sleeping"
  sleep 1
done

echo "PostgreSQL is up - running migrations"
./migration up
if [ $? -ne 0 ]; then
    echo "Migration failed"
    exit 1
fi

echo "Starting the application..."
exec ./engagesync