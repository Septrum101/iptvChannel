FROM golang:alpine AS builder

RUN go env -w GOPROXY=https://goproxy.cn,direct

WORKDIR /build
COPY . .
RUN go mod tidy
RUN go mod download
RUN go build -ldflags="-s -w" -trimpath -v -o /app/main main.go

FROM alpine

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.nju.edu.cn/g' /etc/apk/repositories
RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/main /app/main

CMD ["./main"]