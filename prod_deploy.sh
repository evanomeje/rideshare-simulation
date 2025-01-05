#!/bin/bash

# Build the React app locally
echo "Building the React app..."
cd rideshare-frontend
yarn build
cd ..

# Create a tarball of the build directory
echo "Creating build archive..."
tar -czf frontend-build.tar.gz rideshare-frontend/build

# Copy the build files to the server
echo "Copying build files to server..."
scp frontend-build.tar.gz evan@app.evanomeje.xyz:/home/evan/newserver/

# Extract the build files on the server and clean up
echo "Extracting build files on server..."
ssh evan@app.evanomeje.xyz "cd /home/evan/newserver && \
    tar -xzf frontend-build.tar.gz && \
    rm frontend-build.tar.gz"

# Commit the build files to Git
echo "Committing build files to Git..."
git add .
git commit -m "Build frontend for production"

# Push changes to the remote repository
echo "Pushing changes to Git..."
git push

# Deploy on the server
echo "Deploying to production server..."
ssh -t evan@app.evanomeje.xyz "cd /home/evan/newserver && ./prod_deploy.sh"