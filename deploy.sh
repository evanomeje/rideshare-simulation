#!/bin/bash
# SSH into the production server and run the prod_deploy.sh script
#ssh -t evan@app.evanomeje.xyz "cd /home/evan/newserver && ./prod_deploy.sh"


#!/bin/bash

# Build the React app locally
echo "Building the React app..."
cd rideshare-frontend
yarn build
cd ..

# Commit the build files to Git (optional)
echo "Committing build files to Git..."
git add .
git commit -m "Build frontend for production"

# Push changes to the remote repository
echo "Pushing changes to Git..."
git push

# Step 4: SSH into the production server and run the prod_deploy.sh script
echo "Deploying to production server..."
ssh -t evan@app.evanomeje.xyz "cd /home/evan/newserver && ./prod_deploy.sh"
