FROM golang:1.21-alpine AS build_base
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN cd cmd/terrarium-bot && \
    go build -o /build/out/my-app .

# Start fresh from a smaller image
FROM alpine:3.17.2
RUN apk add ca-certificates
COPY --from=build_base /build/out/my-app /app/terrarium-bot
WORKDIR /app
CMD ["/app/terrarium-bot"]
