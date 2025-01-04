#!/bin/bash
SECONDS=0

# Navigate to your project directory
cd /home/evan/newserver

#cd /home/app
# Function to print messages with a separator
msg () {
  echo -e "$1\n--------------------\n"
}


msg "Pulling from GitHub"
git pull

msg "Building the 'app' image"
sudo docker build --tag app .

msg "Stopping containers"
sudo docker compose down

msg "Starting containers"
sudo docker compose up -d

msg "Pruning stale Docker images"
sudo docker image prune -f

duration=$SECONDS

echo
msg "Deploy finished in $(($duration % 60)) seconds."
msg "Press Enter to exit"
read
