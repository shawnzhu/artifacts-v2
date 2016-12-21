FROM alpine:3.4

EXPOSE 8080

ENV GIN_MODE=release

ADD bin/travis-artifacts /travis-artifacts

ENTRYPOINT ["/travis-artifacts"]
