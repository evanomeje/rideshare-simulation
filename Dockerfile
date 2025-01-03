# Use the official Go image as the base image
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod ./

# Copy the Go sum files
COPY go.sum ./

# Copy postgres stuff
COPY postgres ./postgres

# Download dependencies
RUN go mod download

# Copy the source code
COPY *.go ./
COPY static ./static

# Build the Go binary
RUN go build -o /app/rideshare-simulation

# Expose the port the app will run on
EXPOSE 8080

# Command to run the application
CMD [ "/app/rideshare-simulation" ]