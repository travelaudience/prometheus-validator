package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

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
		Name: "alerts_without_playbook",
		Help: "The alerts that don't have a Playbook link"},
		[]string{
			"alert",
			"owner",
		},
	)
	httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Count of all HTTP requests",
	}, []string{"code", "method"})

	httpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "Duration of all HTTP requests",
	}, []string{"code", "handler", "method"})
)

func apiGet(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil

}

func clear(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

func checkAlerts(url string, noPlaybookAlerts *[]AlertRule) []AlertRule {
	clear(noPlaybookAlerts)
	alertrules := AlertRuleApiResponse{}
	apiResp, err := apiGet(url)
	if err != nil {
		log.Fatalf("Couldn't read apiResp from url %s : %v. \n", url, err)
	}
	err = json.Unmarshal(apiResp, &alertrules)
	if err != nil {
		log.Fatalf("Couldn't marshel the json . %v", err)
	}

	for _, group := range alertrules.Data.Groups {
		for _, rule := range group.Rules {
			noPlaybook := (rule.Type == "alerting") && (rule.Annotations.Playbook == "") && (rule.Annotations.PlayBook == "")
			if noPlaybook {
				*noPlaybookAlerts = append(*noPlaybookAlerts, rule)
			}
		}
	}
	return *noPlaybookAlerts
}

func recordMetrics(queryInterval time.Duration, url string, noPlaybookAlerts *[]AlertRule) {
	checkAlerts(url, noPlaybookAlerts)
	for _, alert := range *noPlaybookAlerts {
		alertsNoPlaybook.With(prometheus.Labels{"alert": alert.Name, "owner": alert.Labels.Owner}).Set(1)
	}
	ticker := time.NewTicker(queryInterval * time.Minute)
	for range ticker.C {
		checkAlerts(url, noPlaybookAlerts)
		for _, alert := range *noPlaybookAlerts {
			alertsNoPlaybook.With(prometheus.Labels{"alert": alert.Name, "owner": alert.Labels.Owner}).Set(1)
		}
	}
}

func main() {
	var noPlaybookAlerts []AlertRule
	fmt.Println(os.Getenv("PROMETHEUS_RULES_API_URL"))
	go recordMetrics(1, os.Getenv("PROMETHEUS_RULES_API_URL"), &noPlaybookAlerts)
	r := prometheus.NewRegistry()
	r.MustRegister(alertsNoPlaybook)
	r.MustRegister(httpRequestsTotal)
	r.MustRegister(httpRequestDuration)

	foundHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from example application."))
	})
	notfoundHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	foundChain := promhttp.InstrumentHandlerDuration(
		httpRequestDuration.MustCurryWith(prometheus.Labels{"handler": "found"}),
		promhttp.InstrumentHandlerCounter(httpRequestsTotal, foundHandler),
	)

	http.Handle("/", foundChain)
	http.Handle("/err", promhttp.InstrumentHandlerCounter(httpRequestsTotal, notfoundHandler))

	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":2112", nil))
}
