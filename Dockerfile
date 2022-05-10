FROM golang:1.18.0 AS builder

ADD . /app
WORKDIR /app
# GOOS/GOARCH as you build not from go alpine
RUN GOOS=linux GOARCH=amd64 go build -o go-postgresql-app ./cmd/go-postgresql-crud-example

FROM alpine:3.15 AS app
WORKDIR /app
COPY --from=builder /app/go-postgresql-app /app
COPY --from=builder /app/cmd/go-postgresql-crud-example/config.yaml /app
CMD ["/app/go-postgresql-app"]