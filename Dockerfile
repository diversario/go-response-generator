FROM golang:1.13.0-alpine3.10

COPY . /app
WORKDIR /app
RUN go build -o app
RUN go build -o client client/client.go

CMD ["./app"]