FROM quay.io/prometheus/busybox:latest
LABEL MAINTAINER="Audrius Karabanovas <auk@trustpilot.com>"

COPY .build/linux-amd64/beat-exporter /bin/beat-exporter

EXPOSE      9479
ENTRYPOINT  [ "/bin/beat-exporter" ]