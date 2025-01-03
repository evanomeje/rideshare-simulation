#!/bin/bash
SECONDS=0

# Navigate to your project directory
cd /home/evan/newserver

# Function to print messages with a separator
msg () {
  echo -e "$1\n--------------------\n"
}

# Stop and remove the running Docker container (if any)
msg "Stopping and removing existing Docker container"
sudo docker stop rideshare-simulation || true  # Stop the container if it exists
sudo docker rm rideshare-simulation || true    # Remove the container if it exists

# Pull the latest code from GitHub
msg "Pulling from GitHub"
git pull

# Build the Docker image
msg "Building Docker image"
sudo docker build -t chima2767/rideshare-simulation:latest .

# Start the Docker container
msg "Starting Docker container"
sudo docker run \
-d \
--name rideshare-simulation \
--expose 443 \
-p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt \
-e SERVER_ENV=PROD \
chima2767/rideshare-simulation:latest

# Prune unused Docker images to free up space
msg "Pruning unused Docker images"
sudo docker image prune -f

# Calculate and display the deployment time
duration=$SECONDS
echo
msg "Deploy finished in $(($duration % 60)) seconds."

# Wait for user input before exiting (useful for monitoring)
msg "Press Enter to exit"
read