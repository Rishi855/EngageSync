#!/bin/bash

echo "Stopping containers on ports 5432 and 8080..."
docker ps -q --filter "publish=5432" | xargs -r docker stop
docker ps -q --filter "publish=8080" | xargs -r docker stop

echo "Pull repo and postgres images from docker-compose..."
docker-compose pull

echo "Starting containers with docker-compose..."
docker-compose up
