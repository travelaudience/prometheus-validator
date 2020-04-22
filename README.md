# Prometheus validator

This repo should contain the code for a prometheus sidecontainer that will allow different validation functionalities:
* Query alerts without playbook and expose them as metrics.

### The app will serve the following metrics under `/metrics` endpoint:
- Go app default Prometheus metrics
- `prometheus_validator_alerts_without_playbook` - Gauge metric with labels for alerts that don't have `playbook` annotation (`play_book`, `runbook` or `runbook_url` are also checked) defined and their owner.

Example:
`prometheus_validator_alerts_without_playbook{alert_name="ClockSkewDetected",alert_owner=""}`

#### The app is built into a docker image and served by
    quay.io/travelaudience/prometheus-validator:[git-tag]

### How to use the docker image?
The image could be added as a Prometheus sidecontainer by using:
1. Prometheus helm chart `server.sidecarContainers`
2. Prometheus Operator helm chart `prometheus.prometheusSpec.containers`

Container config:

```yaml
containers:
  - name: prometheus-validator
    image: quay.io/travelaudience/prometheus-validator:0.1.2
    ports:
      - name: metrics
        containerPort: 8080
    command: ["/util"]
    env:
      - name: PROMETHEUS_URL
        value: "http://localhost:9090"
```
After deploying it as a side-container, you will need to make sure the metrics are being scraped by Prometheus.
This can be done by:
1. ServiceMonitor configuration (for `prometheus-operator`).
2. Prometheus ScrapeConfigs.


### Issue tracker:
https://github.com/travelaudience/prometheus-validator/issues
