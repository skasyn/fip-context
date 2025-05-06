# build stage
FROM golang:alpine as build

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o fip-context

# deploy stage
FROM alpine:latest

WORKDIR /

COPY --from=build /app/fip-context /fip-context
COPY --from=build /app/.env /.env

EXPOSE 8080

ENTRYPOINT ["/fip-context"]