FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o elimulocal .

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/elimulocal .

COPY --from=builder /app/templates ./templates

COPY --from=builder /app/static ./static

RUN mkdir -p uploads

EXPOSE 8080

CMD ["./elimulocal"]

