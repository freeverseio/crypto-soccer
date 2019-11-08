#!/bin/bash

echo "Starting ..."
echo "fdggreTRGDSBw45ergseth4hDGHD" | docker login -u freeversedigitalocean --password-stdin
docker-compose -f docker-compose.development.yml pull
docker-compose -f docker-compose.development.yml down --remove-orphans
docker-compose -f docker-compose.development.yml up -d
docker image prune -f
docker volume prune -f
docker logout

