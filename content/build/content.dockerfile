FROM golang:1.22.2-alpine as builder

RUN mkdir /app

COPY .. /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o content ./cmd

RUN chmod +x /app/content

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/content /app

COPY /build/content-cert.pem /build/content-cert.pem
COPY /build/content-key.pem /build/content-key.pem
COPY /configs/content-config.yaml /configs/content-config.yaml

CMD [ "app/content" ]