FROM golang:1.21.6 AS build

WORKDIR /app

# Set necessary environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Copy Go module files and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o /nhl-mis

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /nhl-mis /nhl-mis

EXPOSE 8080

ENTRYPOINT ["/nhl-mis"]
