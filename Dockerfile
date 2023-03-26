#FROM golang:alpine3.17
FROM golang:1.20

WORKDIR /app
COPY . .
#RUN apt-get update
#RUN apt-get -y install postgresql-client
#RUN chmod +x scripts/wait_for_postgres.sh
RUN go mod download
RUN go build cmd/main.go
EXPOSE 8001
CMD ["./main"]
