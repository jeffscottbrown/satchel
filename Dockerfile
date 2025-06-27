FROM golang:1.24.4-alpine AS appbuilder

RUN apk update && apk add --no-cache build-base go
WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO_ENABLED=1
RUN go build -tags=production -o satchel .

FROM alpine:latest

ENV GIN_MODE=release

ARG PROJECT_ID
ENV PROJECT_ID=$PROJECT_ID

WORKDIR /app

COPY --from=appbuilder /build/satchel ./

CMD ["./satchel"]
