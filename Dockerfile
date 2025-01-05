FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY postgres ./postgres
RUN go mod download

COPY *.go ./

# Create and copy the React build files
COPY rideshare-frontend/build/ /app/rideshare-frontend/build/

# Verify the files are copied (for debugging)
RUN ls -la /app/rideshare-frontend/build

RUN go build -o /app/main

EXPOSE 8080

CMD [ "/app/main" ]