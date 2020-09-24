package jaeger_negroni

import (

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
    cfg.PreRequestMetrics = append(cfg.PreRequestMetrics, metric)
}

// function used to add new pre-request metrics to jaeger configuration
func(cfg JaegerNegroniConfig) AddPostRequestMetric(metric JaegerMetric) {
    cfg.PostRequestMetrics = append(cfg.PostRequestMetrics, metric)
}

