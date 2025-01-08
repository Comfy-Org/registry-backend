# Start from the official Go image to build the binary.
# Use the same version of Go as your project (1.21.5).
FROM golang:1.23.0 AS builder

# Set the working directory inside the container.
WORKDIR /go/src/app

# Copy the Go Modules manifests and download modules to leverage Docker cache.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code.
COPY . .

# Build the Go app.
# -o /go/bin/app: Compile the binary to /go/bin/app.
# You might need to adjust the build command depending on your project's structure.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /go/bin/app

# Use a small base image to create a minimal final image.
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Copy the pre-built binary file from the previous stage.
COPY --from=builder /go/bin/app /go/bin/app

# Run the binary.
ENTRYPOINT ["/go/bin/app"]
