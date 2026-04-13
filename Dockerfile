FROM golang:1.26.2-alpine3.23 AS builder

WORKDIR /app

COPY . .
COPY go.mod ./

RUN go mod tidy

RUN go build app/main.go 

FROM golang:1.26.2-alpine3.23 AS runner
WORKDIR /app

COPY --from=builder /app/main ./

EXPOSE 6379

CMD ["./main"] 