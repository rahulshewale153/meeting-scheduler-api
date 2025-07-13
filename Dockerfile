FROM golang:1.23-alpine

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the source code into the container
COPY . .
# Build the Go application
RUN go build -o meeting-scheduler-api main.go
CMD ["./meeting-scheduler-api"]