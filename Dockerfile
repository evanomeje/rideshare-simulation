# Use the official Go image as the base image
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod ./
COPY go.sum ./

# Copy postgres stuff
COPY postgres ./postgres

# Download dependencies
RUN go mod download

# Copy the source code
COPY *.go ./

# Create the build directory
RUN mkdir -p ./rideshare-frontend/build

# Copy the React build files
COPY rideshare-frontend/build/ ./rideshare-frontend/build/

# Verify the files are copied correctly
RUN ls -la ./rideshare-frontend/build && \
    ls -la ./rideshare-frontend/build/static/js && \
    ls -la ./rideshare-frontend/build/static/css

# Build the Go binary
RUN go build -o /app/main

# Expose the port the app will run on
EXPOSE 8080

# Command to run the application
CMD [ "/app/main" ]