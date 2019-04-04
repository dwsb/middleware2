#!/bin/bash

docker-compose -f docker/docker-compose.yml stop
docker-compose -f docker/docker-compose.yml rm --force
docker-compose -f docker/docker-compose.yml up -d

