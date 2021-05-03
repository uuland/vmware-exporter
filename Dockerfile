FROM golang:1.16 AS builder
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM scratch
COPY --from=builder /go/src/app/vmware-exporter /vmware-exporter
EXPOSE 9512
ENTRYPOINT ["/vmware-exporter"]