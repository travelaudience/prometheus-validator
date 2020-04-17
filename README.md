# Prometheus validator

This repo should contain the code for a prometheus sidecontainer that will allow different validation functionalities:
* Query alerts without playbook and expose them as metrics.

The app will serve the following metrics under `/metrics` endpoint:
- Go app default Prometheus metrics
- `prometheus_validator_alerts_without_playbook` - Gauge metric with labels for alerts that don't have `playbook` or `play_book` defined and their owner.
Example:
`prometheus_validator_alerts_without_playbook{alert_name="ClockSkewDetected",alert_owner=""}`

