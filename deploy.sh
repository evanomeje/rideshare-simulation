#!/bin/bash
SECONDS=0

# Build the React app
echo "Building React app..."
cd rideshare-frontend
yarn build
cd ..

# Sync frontend files
echo "Syncing frontend files to server..."
ssh evan@app.evanomeje.xyz "mkdir -p /home/evan/newserver/rideshare-frontend/build"
scp -r rideshare-frontend/build/* evan@app.evanomeje.xyz:/home/evan/newserver/rideshare-frontend/build/

# Commit and push changes
echo "Committing and pushing changes..."
git add .
git commit -m "Build frontend"
git push

# Deploy to server
echo "Deploying to server..."
ssh evan@app.evanomeje.xyz "cd /home/evan/newserver && \
git pull && \
docker compose down && \
docker build --tag app . && \
docker compose up -d && \
docker logs newserver-app-1"

duration=$SECONDS
echo "Deploy finished in $(($duration % 60)) seconds."