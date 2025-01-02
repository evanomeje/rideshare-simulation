#!/bin/bash
# SSH into the production server and run the prod_deploy.sh script
sshcmd="ssh -t evan@app.evanomeje.xyz"
$sshcmd screen -S "deployment" /home/evan/newserver/prod_deploy.sh

