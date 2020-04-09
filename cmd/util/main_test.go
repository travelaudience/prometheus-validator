package main

import (
	"reflect"
	"testing"
	"time"
)

func Test_recordMetrics(t *testing.T) {
	type args struct {
		queryInterval    time.Duration
		url              string
		noPlaybookAlerts *[]AlertRule
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{name: "invalid-interval",
			args: args{
				queryInterval:    1,
				url:              "http://localhost:8080/api/v1/rules?type=alert",
				noPlaybookAlerts: &[]AlertRule{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// recordMetrics(tt.args.queryInterval, tt.args.url, tt.args.noPlaybookAlerts)
		})
	}
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := apiGet(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("apiGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("apiGet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkAlerts(t *testing.T) {
	type args struct {
		url              string
		noPlaybookAlerts *[]AlertRule
	}
	tests := []struct {
		name string
		args args
		want []AlertRule
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkAlerts(tt.args.url, tt.args.noPlaybookAlerts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkAlerts() = %v, want %v", got, tt.want)
			}
		})
	}
}
