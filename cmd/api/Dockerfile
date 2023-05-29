FROM golang:1.20

COPY /etc/letsencrypt/live/mailbx.ru /app/certs
WORKDIR /app
COPY . .
RUN go mod download
RUN go build cmd/api/main.go

EXPOSE 8001

