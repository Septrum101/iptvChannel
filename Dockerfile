FROM golang:alpine AS builder

RUN go env -w GOPROXY=https://goproxy.cn,direct

WORKDIR /build
COPY . .
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.nju.edu.cn/g' /etc/apk/repositories
RUN apk update --no-cache && apk add --no-cache git
RUN go mod tidy
RUN go build -ldflags="-s -w -X 'github.com/thank243/iptvChannel/config.commit=$(git rev-parse --short HEAD)'" -trimpath -v -o /app/main main.go

FROM alpine

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.nju.edu.cn/g' /etc/apk/repositories
RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/main /app/main

CMD ["./main"]