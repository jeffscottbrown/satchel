FROM golang:1.24.3-alpine AS appbuilder

RUN apk update && apk add --no-cache build-base go
WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY main.go .
COPY server/ server/

ENV CGO_ENABLED=1
RUN go build -o satchel .

FROM alpine:latest

ENV GIN_MODE=release

WORKDIR /app

COPY --from=appbuilder /build/satchel ./

CMD ["./satchel"]
