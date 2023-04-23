FROM golang:alpine3.17

WORKDIR /app
COPY . .
RUN go mod download
RUN go build cmd/api/main.go

EXPOSE 8001
CMD ["./main"]
