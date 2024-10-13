FROM golang:alpine as builder

# Set working directory inside the container
WORKDIR /app

# Copy the Go application source code
COPY . .

# Build the Go binary statically
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o weiba

# Step 2: Create a lightweight final image using Alpine
FROM alpine:latest

# Set working directory in the final image
WORKDIR /app

# Copy the statically compiled Go binary from the builder stage
COPY --from=builder /app/weiba .
COPY ./data.json /app/data.json

# Command to run the application
CMD ["./weiba"]