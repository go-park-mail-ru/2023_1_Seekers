FROM golang:alpine3.17

WORKDIR /app
COPY . .
RUN go mod download
RUN go build cmd/file_storage/main.go

EXPOSE 8005
CMD ["./main"]
