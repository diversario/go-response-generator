FROM golang:1.15-alpine

COPY . /app
WORKDIR /app
RUN go build -o app
RUN go build -o client client/client.go

CMD ["./app"]
