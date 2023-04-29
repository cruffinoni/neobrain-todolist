FROM golang:1.20

WORKDIR /app

COPY . .
RUN go mod vendor
RUN go build -o main cmd/api/main.go

EXPOSE 8080

CMD ["./main"]
