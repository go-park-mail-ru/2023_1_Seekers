FROM golang:1.18 AS builder

WORKDIR /app
COPY . .
RUN go build cmd/main.go

FROM golang:1.18

WORKDIR /app

COPY --from=builder /app/main /app/main

EXPOSE 8001
CMD ["./main"]