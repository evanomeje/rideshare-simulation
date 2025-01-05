#!/bin/bash
SECONDS=0

# Build the React app
echo "Building React app..."
cd rideshare-frontend
yarn build
cd ..

# Commit and push changes
echo "Committing and pushing changes..."
git add .
git commit -m "Build frontend"
git push

# Deploy to server
echo "Deploying to server..."
ssh evan@app.evanomeje.xyz "cd /home/evan/newserver && git pull && ./prod_deploy.sh"

duration=$SECONDS
echo "Deploy finished in $(($duration % 60)) seconds."