#!/bin/bash

echo "Stopping containers on ports 5432 and 8080..."
docker ps -q --filter "publish=5432" | xargs -r docker stop
docker ps -q --filter "publish=8080" | xargs -r docker stop

echo "Removing containers on ports 5432 and 8080..."
docker ps -a -q --filter "publish=5432" | xargs -r docker rm
docker ps -a -q --filter "publish=8080" | xargs -r docker rm

echo "Pulling latest images..."
docker pull postgres:15
docker pull rushikesh855/engagesync-backend

echo "Creating Docker network 'engagesync-net' (if not exists)..."
docker network inspect engagesync-net >/dev/null 2>&1 || docker network create engagesync-net

echo "Starting Postgres container..."
docker run -d --name engagesync-db \
  --network engagesync-net \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=root \
  -e POSTGRES_DB=engagesyncdb \
  -v engagesync_postgres-data:/var/lib/postgresql/data \
  -p 5432:5432 \
  postgres:15

echo "Starting Backend container..."
docker run -d --name engagesync-backend \
  --network engagesync-net \
  -e DB_USER=postgres \
  -e DB_PASSWORD=root \
  -e DB_NAME=engagesyncdb \
  -e DB_HOST=engagesync-db \
  -e DB_PORT=5432 \
  -e DB_SSLMODE=disable \
  -p 8080:8080 \
  rushikesh855/engagesync-backend

echo "All done! Containers are running."
