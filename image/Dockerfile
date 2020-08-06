FROM alpine:latest as builder
ADD https://dl.bintray.com/flant/werf/v1.1.21+fix20/werf-linux-amd64-v1.1.21+fix20 /werf
RUN chmod +x /werf
RUN apk add git
ENTRYPOINT ["/werf"]
