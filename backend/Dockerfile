FROM golang:1.19-alpine3.16 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify


COPY . .
RUN go build -o backend

FROM alpine:3.16
EXPOSE 80

WORKDIR /app

COPY --from=builder /build/backend /app/backend

ENTRYPOINT [ "/app/backend", "--port", "80" ]