FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY iptvChannel /app/iptvChannel
ENTRYPOINT ["/app/iptvChannel"]