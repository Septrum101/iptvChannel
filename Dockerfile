FROM golang:alpine AS builder

ARG VERSION

WORKDIR /build
COPY . .

RUN go mod tidy
RUN go build -trimpath -ldflags="-s -w \
    -X 'github.com/thank243/iptvChannel/config.date=$(date -Iseconds)' \
    -X 'github.com/thank243/iptvChannel/config.version=$VERSION' \
    " -v -o iptvChannel main.go

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /build/iptvChannel /app/iptvChannel
ENTRYPOINT ["/app/iptvChannel"]