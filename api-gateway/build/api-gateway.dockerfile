FROM golang:1.22.2-alpine as builder

RUN mkdir /app

COPY .. /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o api-gateway ./cmd

RUN chmod +x /app/api-gateway

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/api-gateway /app

COPY /build/api-gateway-cert.pem /build/api-gateway-cert.pem
COPY /build/api-gateway-key.pem /build/api-gateway-key.pem
COPY /configs/api-gateway-config.yaml /configs/api-gateway-config.yaml

CMD [ "app/api-gateway" ]