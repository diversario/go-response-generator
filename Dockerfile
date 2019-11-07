FROM golang:1.13.0-alpine3.10

COPY . /app
WORKDIR /app
RUN go build -o server

CMD ["./server"]