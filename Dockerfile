FROM alpine:3.4

EXPOSE 443

ENV GIN_MODE=release

RUN apk add --no-cache ca-certificates

ADD bin/travis-artifacts /travis-artifacts

ENTRYPOINT ["/travis-artifacts"]
