
FROM alpine:latest

RUN apk update && apk add --no-cache libc6-compat
COPY ./product /go/src/
COPY ./config.yaml /

CMD ["/go/src/product"]