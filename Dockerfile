FROM golang:1.21.6-alpine3.19 AS builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o main main.go

FROM alpine:3.19
WORKDIR /app
COPY .env .
COPY --from=builder /app/main .

EXPOSE 8080
CMD [ "/app/main" ]