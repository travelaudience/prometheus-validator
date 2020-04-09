FROM busybox:1.31.1

ADD ./bin/util /

# ENV PROMETHEUS_RULES_API_URL="http://localhost:9090/api/v1/rules?type=alert"
ENV PROMETHEUS_RULES_API_URL="http://prometheus-operated.monitoring.svc.cluster.local:9090/api/v1/rules?type=alert"

CMD ["/util"]

