#!/bin/sh

sh ./setup.sh

docker compose up -d

cd backend && go run . -command server
