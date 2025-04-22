FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o snowflake ./cmd/snowflake.go

FROM scratch
WORKDIR /app/

COPY --from=builder /app/snowflake .

EXPOSE 8080
EXPOSE 5051

CMD ["./snowflake"]
