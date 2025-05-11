# Use Go base image
FROM golang:1.22-alpine

# Set working directory
WORKDIR /app

# Install required packages
RUN apk add --no-cache git

# Copy Go files
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build Go app
RUN go build -o engagesync .

# Run the application
CMD ["./engagesync"]
