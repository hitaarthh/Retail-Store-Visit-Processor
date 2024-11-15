# Use the official Golang image as a build stage
FROM golang:1.19 AS builder

# Set the working directory in the container
WORKDIR /app

# Copy Go modules and dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o main .

# Use a minimal image to run the app
FROM gcr.io/distroless/base-debian11

# Set the working directory for the final image
WORKDIR /app

# Copy the compiled binary from the builder
COPY --from=builder /app/main .

# Copy additional files like CSV if needed
COPY StoreMasterAssignment.csv .

# Copy the frontend folder
COPY frontend ./frontend

# Expose the port
EXPOSE 8081

# Run the binary
CMD ["/app/main"]
