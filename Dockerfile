FROM quay.io/prometheus/golang-builder as builder

ADD .   /go/src/github.com/trustpilot/beat-exporter
WORKDIR /go/src/github.com/trustpilot/beat-exporter

RUN make

FROM        quay.io/prometheus/busybox:latest
MAINTAINER  Audrius Karabanovas <auk@trustpilot.com>

COPY --from=builder /go/src/github.com/trustpilot/beat-exporter/beat-exporter  /bin/beat-exporter

EXPOSE      9479
ENTRYPOINT  [ "/bin/beat-exporter" ]