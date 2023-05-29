FROM golang:1.20

COPY /etc/letsencrypt/live/mailbx.ru /app/certs

WORKDIR /app
COPY . .
RUN go mod download

EXPOSE 8001
