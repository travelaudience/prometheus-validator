package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"testing"
)

func apiResp(file string) []byte {
	r, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("ReadFile failed, %v", err)
	}
	fmt.Printf("r is: %s", r)
	return r
}

func Test_apiGet(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "invalid-url",
			args: args{
				url: "12hgfhfhf",
			},
			want:    []byte{},
			wantErr: true,
		},
		{
			name: "non-exist-url",
			args: args{
				url: "http://www.12hgfhfhf.com",
			},
			want:    []byte{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &PrometheusClient{
				&http.Client{},
				tt.args.url,
			}
			got, err := client.apiGet()
			if (err != nil) != tt.wantErr {
				t.Errorf("apiGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("apiGet() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func Test_checkAlerts(t *testing.T) {
	type args struct {
		apiResp          []byte
		noPlaybookAlerts *[]AlertRule
	}

	tests := []struct {
		name    string
		args    args
		want    *[]AlertRule
		wantErr bool
	}{
		{
			name: "playbook-exist",
			args: args{
				apiResp:          apiResp("../tests-data/playbook.txt"),
				noPlaybookAlerts: &[]AlertRule{},
			},
			want:    &[]AlertRule{},
			wantErr: false,
		},
		{
			name: "playbook-exist-non-clear-struct",
			args: args{
				apiResp: apiResp("../tests-data/playbook.txt"),
				noPlaybookAlerts: &[]AlertRule{
					{
						Name:     "test",
						Query:    "sum(something)",
						Duration: 20,
					},
				},
			},
			want:    &[]AlertRule{},
			wantErr: false,
		},
		{
			name: "play_book-exist",
			args: args{
				apiResp:          apiResp("../tests-data/play_book.txt"),
				noPlaybookAlerts: &[]AlertRule{},
			},
			want:    &[]AlertRule{},
			wantErr: false,
		},
		{
			name: "multiple-alerts-rules",
			args: args{
				apiResp:          apiResp("../tests-data/multiple-rules.txt"),
				noPlaybookAlerts: &[]AlertRule{},
			},
			want:    &[]AlertRule{},
			wantErr: false,
		},
		{
			name: "empty-json",
			args: args{
				apiResp:          apiResp("../tests-data/empty-json.txt"),
				noPlaybookAlerts: &[]AlertRule{},
			},
			want:    &[]AlertRule{},
			wantErr: false,
		},
		{
			name: "unvalid-field-type",
			args: args{
				apiResp:          apiResp("../tests-data/unvalid-field-type.txt"),
				noPlaybookAlerts: &[]AlertRule{},
			},
			want:    &[]AlertRule{},
			wantErr: true,
		},
		{
			name: "invalid-json",
			args: args{
				apiResp:          apiResp("../tests-data/invalid-json.txt"),
				noPlaybookAlerts: &[]AlertRule{},
			},
			want:    &[]AlertRule{},
			wantErr: true,
		},
		{
			name: "no-playbook",
			args: args{
				apiResp:          apiResp("../tests-data/no_playbook.txt"),
				noPlaybookAlerts: &[]AlertRule{},
			},
			want: &[]AlertRule{
				{
					Name:     "TargetUp",
					Query:    "query3",
					Duration: 600,
					Labels: struct {
						Owner    string "json:\"owner,omitempty\""
						Severity string "json:\"severity,omitempty\""
					}{
						Severity: "warning",
					},
					Annotations: struct {
						Description string "json:\"description,omitempty\""
						Playbook    string "json:\"playbook,omitempty\""
						PlayBook    string "json:\"play_book,omitempty\""
						Summary     string "json:\"summary,omitempty\""
					}{},
					Alerts: []Alert{},
					Health: "ok",
					Type:   "alerting",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(*tt.want) == 0 {
				clearObjectValues(tt.want)
			}
			err := checkAlerts(tt.args.apiResp, tt.args.noPlaybookAlerts)
			fmt.Printf("error is %v", err)
			if !reflect.DeepEqual(*tt.args.noPlaybookAlerts, *tt.want) || (err != nil && !tt.wantErr) {
				t.Errorf("Test failed %s : noPlaybookAlert value: %v , noPlaybookAlert type: %T, want value: %v , want type: %T", tt.name, *tt.args.noPlaybookAlerts, *tt.args.noPlaybookAlerts, *tt.want, *tt.want)
			}
		})
	}
}
