#!/bin/bash

echo "🛠 Running database migrations..."
migrate -path /app/migrations -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=$DB_SSLMODE" up

echo "🚀 Starting EngageSync app..."
./engagesyncdb