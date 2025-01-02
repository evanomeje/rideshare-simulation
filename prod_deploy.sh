#!/bin/bash
SECONDS=0

# Navigate to your project directory
#cd ~/Another/rideshare-simulation
cd /home/evan/newserver

# Function to print messages with a separator
msg () {
  echo -e "$1\n--------------------\n"
}

# Stop the running app (if any)
msg "Stopping app"
sudo pkill app || true  # The `|| true` ensures the script doesn't fail if no app is running

# Pull the latest code from GitHub
msg "Pulling from GitHub"
git pull

# Build the Go binary
msg "Building Go binary"
go build

# Start the server in the background
msg "Starting server"
nohup sudo ./app &>/dev/null &  # `nohup` keeps the server running after logout

# Calculate and display the deployment time
duration=$SECONDS
echo
msg "Deploy finished in $(($duration % 60)) seconds."

# Wait for user input before exiting (useful for monitoring)
msg "Press Enter to exit"
read