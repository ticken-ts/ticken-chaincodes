FROM golang:1.18-alpine

# Set the Current Working Directory inside the container
WORKDIR /chaincode/ticken-event

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./out/ticken-event .

# This container exposes port 9999 to the outside world
EXPOSE 9999

# Run the binary program produced by `go install`
CMD ["./out/ticken-event"]