FROM golang:1.20-alpine

RUN mkdir /app 

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o server ./cmd/main.go

EXPOSE 8080 8080

CMD ["./server"]