FROM golang:1.22.2-alpine as builder

RUN mkdir /app

COPY .. /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o file-server ./cmd

RUN chmod +x /app/file-server

FROM alpine:latest

RUN apk update && apk add --no-cache ffmpeg

RUN mkdir /app

COPY --from=builder /app/file-server /app

COPY /configs/file-server-config.yaml /configs/file-server-config.yaml

CMD [ "app/file-server" ]