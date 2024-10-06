# Use the official Golang image to create a build artifact.
# This is the first stage of a multi-stage build.
FROM golang:1.23.1-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o xenigo .

# Start a new stage from scratch
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/xenigo .

# Copy and rename the sample file to config.yaml
# COPY example.yaml config.yaml

# [Placeholder] Expose port 3333 to the outside world
# EXPOSE 3334

# Command to run the executable
CMD ["./xenigo"]