# syntax=docker/dockerfile:1
FROM golang:1.20.2 as builder

WORKDIR /usr/src/app

COPY go.mod ./
COPY go.sum ./
COPY Makefile ./
COPY . .

RUN go mod download

# -ldflags "-s -w" disables some debug related functionality
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -x -v -o /app ./cmd/main.go

# Pulling image for the service to run on
FROM alpine:3.14

WORKDIR /usr/src/app

COPY --from=builder /app /app
COPY --from=builder /usr/src/app/assets /usr/src/app/assets

# Set the command to run the application
CMD ["/app"]