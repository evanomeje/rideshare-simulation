#!/bin/bash
# SSH into the production server and run the prod_deploy.sh script
ssh -t evan@app.evanomeje.xyz "cd /home/evan/newserver && ./prod_deploy.sh"
