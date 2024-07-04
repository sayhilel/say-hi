# Start from the latest golang base image
FROM golang:latest

WORKDIR /app

RUN git clone https://github.com/sayhilel/say-hi.git .

RUN go mod tidy

RUN go build -o main cmd/main.go

EXPOSE 8080

CMD ["./main"]

