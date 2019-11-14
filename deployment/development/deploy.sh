#!/bin/bash

echo "Starting ..."
echo "fdggreTRGDSBw45ergseth4hDGHD" | docker login -u freeversedigitalocean --password-stdin
docker-compose pull
docker-compose down --remove-orphans
docker-compose up -d
docker image prune -f
docker volume prune -f
docker logout

