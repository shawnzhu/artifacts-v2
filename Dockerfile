FROM alpine:3.4

EXPOSE 443

ENV GIN_MODE=release

ADD bin/travis-artifacts /travis-artifacts

ENTRYPOINT ["/travis-artifacts"]
