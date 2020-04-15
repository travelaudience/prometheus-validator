FROM golang:1.14 AS build
WORKDIR  /go/src/github.com/travelaudience/prometheus-validator

COPY ./util ./util

RUN apt-get update && apt-get install -y upx && apt-get clean
RUN go get ./util
RUN CGO_ENABLED=0 GOOS=linux  go build -o util ./util

RUN upx --best util/util


FROM busybox:1.31.1

EXPOSE 8080


COPY --from=build /go/src/github.com/travelaudience/prometheus-validator/util/util /

# ENV PROMETHEUS_RULES_API_URL="http://localhost:9090/api/v1/rules?type=alert"
ENV PROMETHEUS_RULES_API_URL="http://prometheus-operated.monitoring.svc.cluster.local:9090/api/v1/rules?type=alert"

CMD ["/util"]

