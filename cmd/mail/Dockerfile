FROM golang:1.20

WORKDIR /app
COPY . .
RUN go mod download
RUN go build cmd/mail/main.go

EXPOSE 8008

