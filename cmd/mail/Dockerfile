FROM golang:alpine3.17

WORKDIR /app
COPY . .
RUN go mod download
RUN go build cmd/mail/main.go

EXPOSE 8008
CMD ["./main"]
