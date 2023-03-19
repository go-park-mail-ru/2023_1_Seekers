FROM golang:alpine3.17

WORKDIR /app
COPY . .
RUN go build cmd/main.go
EXPOSE 8001
CMD ["./main"]
