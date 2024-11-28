FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o api cmd/api/main.go
RUN go build -o consumer cmd/consumer/main.go

#--------------------------------------------------
FROM gcr.io/distroless/base-debian11 AS api

WORKDIR /app

COPY --from=builder /app/api .

EXPOSE 8080

CMD ["/app/api"]

#--------------------------------------------------
FROM gcr.io/distroless/base-debian11 AS consumer

WORKDIR /app

COPY --from=builder /app/consumer .

CMD ["/app/consumer"]