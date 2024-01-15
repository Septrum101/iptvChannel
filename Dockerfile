FROM golang:alpine AS builder

WORKDIR /build
COPY . .

RUN apk update --no-cache && apk add --no-cache git
RUN go mod tidy
RUN go build -ldflags="-s -w \
     -X 'github.com/thank243/iptvChannel/config.commit=$(git rev-parse --short HEAD) build: $(date)'" \
    -trimpath -v -o /app/main main.go

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/main /app/main

CMD ["./main"]