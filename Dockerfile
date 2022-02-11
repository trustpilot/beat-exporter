FROM golang:1.16.10-alpine3.13 AS builder

WORKDIR /root/go/src/beat-exporter

ENV GOPROXY https://goproxy.cn
ENV GOPATH /root/go
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -o beat-exporter

FROM alpine:3.13 AS final

WORKDIR /bin
COPY --from=builder /root/go/src/beat-exporter/beat-exporter /bin/

EXPOSE 9479
ENTRYPOINT ["/bin/beat-exporter"]