package jaeger_negroni

import (
    "fmt"
    log "github.com/sirupsen/logrus"
)

type JaegerConfig struct {
    JaegerHost  string `json:"jaeger_host"`
    JaegerPort  int    `json:"jaeger_port"`
    ServiceName string `json:"service_name"`
}

type JaegerNegroniConfig struct {
    JaegerConf         JaegerConfig   `json:"jaeger_conf"`
    PreRequestMetrics  []JaegerMetric `json:"pre_request_metrics"`
    PostRequestMetrics []JaegerMetric `json:"post_request_metrics"`
}

// function used to add new pre-request metrics to jaeger configuration
func(cfg JaegerNegroniConfig) AddPreRequestMetric(metric JaegerMetric) {
    log.Debug(fmt.Sprintf("adding new pre-request metric %s", metric))
    cfg.PreRequestMetrics = append(cfg.PreRequestMetrics, metric)
}

// function used to add new pre-request metrics to jaeger configuration
func(cfg JaegerNegroniConfig) AddPostRequestMetric(metric JaegerMetric) {
    log.Debug(fmt.Sprintf("adding new post-request metric %s", metric))
    cfg.PostRequestMetrics = append(cfg.PostRequestMetrics, metric)
}

