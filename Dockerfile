FROM docker.io/alpine:3.14

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY fakenews-server /fakenews-server
ENTRYPOINT ["/fakenews-server"]
