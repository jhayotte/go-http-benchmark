FROM golang:1.21

WORKDIR /app
COPY . .

RUN go mod init beego && go mod tidy && go build -o server

EXPOSE 8080
CMD ["./server"]
