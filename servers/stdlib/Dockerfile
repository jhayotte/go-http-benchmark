FROM golang:1.21

WORKDIR /app
COPY . .

RUN go mod init server_stdlib && go mod tidy && go build -o server

EXPOSE 8080
CMD ["./server"]
