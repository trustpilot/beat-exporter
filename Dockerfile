FROM quay.io/prometheus/busybox:uclibc

COPY beat-exporter /usr/local/bin/beat-exporter

EXPOSE      9479
ENTRYPOINT  [ "beat-exporter" ]
