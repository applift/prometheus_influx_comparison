FROM golang:1.9

RUN mkdir -p /go/src/github.com/gburanov/prometheus_influx_comparison
WORKDIR /go/src/github.com/gburanov/prometheus_influx_comparison

COPY . /go/src/github.com/gburanov/prometheus_influx_comparison

RUN go get -u github.com/golang/dep/cmd/dep
RUN godep restore
RUN rm -f /go/src/github.com/gburanov/prometheus_influx_comparison/prometheus_influx_comparison
RUN go build .

CMD ["/go/src/github.com/gburanov/prometheus_influx_comparison/prometheus_influx_comparison"]
