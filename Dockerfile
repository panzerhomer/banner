FROM golang:alpine

# WORKDIR /app
# ADD go.mod .
# COPY . .
# RUN go build -o main main.go

# ENTRYPOINT ["./main"]

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main ./cmd/main.go
EXPOSE 5000

ENTRYPOINT ["./main"]