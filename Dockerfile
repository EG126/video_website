FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o video_website main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/video_website .

COPY --from=builder /app/config /root/config

COPY --from=builder /app/static /root/static

EXPOSE 8888

CMD ["./video_website"]