FROM golang:1.23-alpine

WORKDIR /app

# Copy dependencies first
COPY go.mod go.sum ./

# Copy source code
COPY server ./server
COPY db ./db

# Build the application
RUN cd server && go build -o main .

EXPOSE 8080

CMD ["./server/main"]