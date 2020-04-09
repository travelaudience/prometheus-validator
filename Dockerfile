FROM golang:1.14 AS build
WORKDIR  /go/src/github.com/travelaudience/prometheus-validator

COPY ./util ./util

RUN go get ./util
RUN CGO_ENABLED=0 GOOS=linux  go build -o util ./util


FROM busybox:1.31.1


COPY --from=build /go/src/github.com/travelaudience/prometheus-validator/util /

# ENV PROMETHEUS_RULES_API_URL="http://localhost:9090/api/v1/rules?type=alert"
ENV PROMETHEUS_RULES_API_URL="http://prometheus-operated.monitoring.svc.cluster.local:9090/api/v1/rules?type=alert"

CMD ["/util"]

