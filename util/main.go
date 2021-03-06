package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusClient struct {
	httpclient *http.Client
	url        string
}

type Alert struct {
	Labels struct {
		AlertName string `json:"alertname,omitempty"`
		Job       string `json:"job,omitempty"`
		Namespace string `json:"namespace,omitempty"`
		Service   string `json:"service,omitempty"`
		Severity  string `json:"severity,omitempty"`
	} `json:"labels"`
	Annotations struct {
		Message string `json:"message,omitempty"`
	} `json:"annotations,omitempty"`
	State    string `json:"state,omitempty"`
	ActiveAt string `json:"activeAt,omitempty"`
	Value    string `json:"value,omitempty"`
}

type AlertRule struct {
	Name     string `json:"name"`
	Query    string `json:"query"`
	Duration int    `json:"duration,omitempty"`
	Labels   struct {
		Owner    string `json:"owner,omitempty"`
		Severity string `json:"severity,omitempty"`
	} `json:"labels"`
	Annotations struct {
		Description string `json:"description,omitempty"`
		Playbook    string `json:"playbook,omitempty"`
		PlayBook    string `json:"play_book,omitempty"`
		RunbookURL  string `json:"runbook_url,omitempty"`
		Runbook     string `json:"runbook,omitempty"`
		Summary     string `json:"summary,omitempty"`
	} `json:"annotations,omitempty"`
	Alerts []Alert `json:"alerts"`
	Health string  `json:"health"`
	Type   string  `json:"type"`
}

type AlertRuleGroup struct {
	Name     string      `json:"name"`
	File     string      `json:"file"`
	Rules    []AlertRule `json:"rules"`
	Interval int         `json:"interval,omitempty"`
}
type AlertRuleApiResponse struct {
	Status string `json:"status"`
	Data   struct {
		Groups []AlertRuleGroup `json:"groups"`
	} `json:"data"`
}

var (
	alertsNoPlaybook = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "prometheus_validator_alerts_without_playbook",
		Help: "The alerts that don't have a Playbook link"},
		[]string{
			"alert_name",
			"alert_owner",
		},
	)
)

func (client *PrometheusClient) apiGet() ([]byte, error) {
	req, err := http.NewRequest("GET", client.url, nil)

	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.httpclient.Do(req)

	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return []byte{}, err
	}

	return body, nil

}

func clearObjectValues(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

func checkAlerts(apiResp []byte, noPlaybookAlerts *[]AlertRule) error {
	clearObjectValues(noPlaybookAlerts)
	alertrules := AlertRuleApiResponse{}
	err := json.Unmarshal(apiResp, &alertrules)
	if err != nil {
		return err
	}

	for _, group := range alertrules.Data.Groups {
		for _, rule := range group.Rules {
			noPlaybook := (rule.Type == "alerting") && (rule.Annotations.Playbook == "") && (rule.Annotations.PlayBook == "") && (rule.Annotations.Runbook == "") && (rule.Annotations.RunbookURL == "")
			if noPlaybook {
				*noPlaybookAlerts = append(*noPlaybookAlerts, rule)
			}
		}
	}
	return nil
}

func recordMetrics(queryInterval time.Duration, noPlaybookAlerts *[]AlertRule, client *PrometheusClient) {
	apiResp, err := client.apiGet()
	if err != nil {
		log.Fatalf("Couldn't read apiResp from url %s : %v. \n", client.url, err)
	}
	err = checkAlerts(apiResp, noPlaybookAlerts)
	if err != nil {
		log.Fatalf("Couldn't marshel the json . %v", err)
	}
	for _, alert := range *noPlaybookAlerts {
		alertsNoPlaybook.With(prometheus.Labels{"alert_name": alert.Name, "alert_owner": alert.Labels.Owner}).Set(1)
	}
	ticker := time.NewTicker(queryInterval * time.Minute)
	for range ticker.C {
		apiResp, err := client.apiGet()
		if err != nil {
			log.Fatalf("Couldn't read apiResp from url %s : %v. \n", client.url, err)
		}
		err = checkAlerts(apiResp, noPlaybookAlerts)
		if err != nil {
			log.Fatalf("Couldn't marshel the json . %v", err)
		}
		for _, alert := range *noPlaybookAlerts {
			alertsNoPlaybook.With(prometheus.Labels{"alert_name": alert.Name, "alert_owner": alert.Labels.Owner}).Set(1)
		}
	}
}

func main() {
	var noPlaybookAlerts []AlertRule
	url := os.Getenv("PROMETHEUS_URL") + "/api/v1/rules?type=alert"
	client := &PrometheusClient{
		&http.Client{},
		url,
	}
	go recordMetrics(1, &noPlaybookAlerts, client)

	prometheus.MustRegister(alertsNoPlaybook)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
