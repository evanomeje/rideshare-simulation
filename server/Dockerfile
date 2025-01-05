FROM golang:1.21-alpine

WORKDIR /app

# Copy dependencies first
COPY go.mod go.sum ./

# Copy configuration
COPY dbconfig.json ./dbconfig.json
COPY .env ./.env

# Copy source code
COPY server ./server
COPY db ./db

# Build the application
RUN cd server && \
    CGO_ENABLED=0 GOOS=linux go build -o main

# Move the binary to the expected location
RUN mv server/main ./main

EXPOSE 8080

CMD ["./main"]
