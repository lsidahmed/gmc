FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o crudApp ./cmd/main.go

EXPOSE 8090

CMD ["./crudApp"]