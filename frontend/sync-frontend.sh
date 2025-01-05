#!/bin/bash
echo "Syncing frontend build to server..."
ssh evan@app.evanomeje.xyz "mkdir -p /home/evan/newserver/rideshare-frontend/build"
scp -r build/* evan@app.evanomeje.xyz:/home/evan/newserver/rideshare-frontend/build/
