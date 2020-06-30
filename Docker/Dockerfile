# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:rc-alpine3.12 as builder
RUN apk add --no-cache --virtual .build-deps gcc musl-dev openssl openssh git 
WORKDIR /go/src/
RUN git clone https://github.com/cheetahfox/nfCollector.git
RUN go get github.com/landoop/tableprinter
RUN go get github.com/spf13/viper
RUN go get github.com/fln/nf9packet
RUN go get github.com/tehmaze/netflow
RUN go get github.com/tehmaze/netflow/ipfix
RUN go get github.com/tehmaze/netflow/netflow1
RUN go get github.com/tehmaze/netflow/netflow5
RUN go get github.com/tehmaze/netflow/netflow6
RUN go get github.com/tehmaze/netflow/netflow7
RUN go get github.com/tehmaze/netflow/netflow9
RUN go get github.com/tehmaze/netflow/session
RUN go get github.com/influxdata/influxdb1-client
WORKDIR /go/src/nfCollector/cmd/nfc
RUN go build
WORKDIR /go/src/nfCollector/cmd/nfc-dump
RUN go build
RUN strip nfc-dump 

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir /etc/nfcol/
RUN mkdir /tmp/nfcol/
COPY --from=builder /go/src/nfCollector/cmd/nfc/nfc .
COPY --from=builder /go/src/nfCollector/cmd/nfc-dump/nfc-dump .
CMD ./nfc
