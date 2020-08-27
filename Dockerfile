# Build
FROM golang:alpine as builder

LABEL maintainer="Andy Jones <andy@andyjones.co>"

RUN apk update && apk add --no-cache git gcc

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o server ./cmd/server
RUN go build -o migrate ./cmd/migrate

# Run
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /usr/src/app/server .
COPY --from=builder /usr/src/app/migrate .
COPY --from=builder /usr/src/app/migrations ./migrations

EXPOSE 8000

ENTRYPOINT ["/root/server"]