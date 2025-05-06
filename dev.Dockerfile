FROM golang:alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o fip-context-dev

EXPOSE 8080
ENTRYPOINT ["/app/fip-context-dev"]