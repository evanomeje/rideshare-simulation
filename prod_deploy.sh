#!/bin/bash
SECONDS=0

# Navigate to your project directory
cd /home/evan/newserver

# Function to print messages with a separator
msg () {
  echo -e "$1\n--------------------\n"
}

# Stop and remove existing containers
msg "Stopping containers"
sudo docker compose down

# Pull the latest code from GitHub
msg "Pulling from GitHub"
git pull

# Start the containers
msg "Starting containers"
sudo docker compose up -d

# Prune unused Docker images
msg "Pruning unused Docker images"
sudo docker image prune -f

# Calculate and display the deployment time
duration=$SECONDS
echo
msg "Deploy finished in $(($duration % 60)) seconds."

# Wait for user input before exiting (useful for monitoring)
msg "Press Enter to exit"
read