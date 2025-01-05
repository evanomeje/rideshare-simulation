#!/bin/bash
echo "Syncing frontend build to server..."
rsync -avz build/ evan@app.evanomeje.xyz:/home/evan/newserver/rideshare-frontend/build/
