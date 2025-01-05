#!/bin/bash
SECONDS=0

# Build the React app
echo "Building React app..."
cd rideshare-frontend
yarn build
cd ..

# Create sync script
cat > sync-frontend.sh << 'EOF'
#!/bin/bash
echo "Syncing frontend build to server..."
ssh evan@app.evanomeje.xyz "mkdir -p /home/evan/newserver/rideshare-frontend/build"
scp -r rideshare-frontend/build/* evan@app.evanomeje.xyz:/home/evan/newserver/rideshare-frontend/build/
EOF

chmod +x sync-frontend.sh

# Sync frontend files
echo "Syncing frontend files to server..."
./sync-frontend.sh

# Commit and push changes
echo "Committing and pushing changes..."
git add .
git commit -m "Build frontend"
git push

# Deploy to server
echo "Deploying to server..."
ssh -t evan@app.evanomeje.xyz "cd /home/evan/newserver && ./prod_deploy.sh"

duration=$SECONDS
echo "Deploy finished in $(($duration % 60)) seconds."