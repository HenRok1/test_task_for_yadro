FROM golang:latest

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN make build

CMD ["./main", "files/test2.txt"]