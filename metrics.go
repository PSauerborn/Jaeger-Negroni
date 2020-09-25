package jaeger_negroni

import (
    "github.com/gin-gonic/gin"
)

var (
    preRequestMetrics = []JaegerMetric{HTTPUrlMetric{}, HTTPRequestMethodMetric{}}
    postRequestMetrics = []JaegerMetric{HTTPResponseStatusMetric{}, HTTPErrorStatusMetric{}}
)

// interface used to add various jaeger metrics to the tracer config
type JaegerMetric interface {
    MetricName() string
    EvaluateMetric(ctx *gin.Context) interface{}
}

// define metric used to add request route to jaeger metrics
type HTTPUrlMetric struct{}

func(metric HTTPUrlMetric) MetricName() string {
    return "http.url"
}

func(metric HTTPUrlMetric) EvaluateMetric(ctx *gin.Context) interface{} {
    return ctx.FullPath()
}

// define metric used to add request method to jaeger metrics
type HTTPRequestMethodMetric struct{}

func(metric HTTPRequestMethodMetric) MetricName() string {
    return "http.method"
}

func(metric HTTPRequestMethodMetric) EvaluateMetric(ctx *gin.Context) interface{} {
    return ctx.Request.Method
}

// define metric used to add request method to jaeger metrics
type HTTPResponseStatusMetric struct{}

func(metric HTTPResponseStatusMetric) MetricName() string {
    return "http.status"
}

func(metric HTTPResponseStatusMetric) EvaluateMetric(ctx *gin.Context) interface{} {
    return ctx.Writer.Status()
}

// define metric used to add request method to jaeger metrics
type HTTPErrorStatusMetric struct{}

func(metric HTTPErrorStatusMetric) MetricName() string {
    return "error"
}

func(metric HTTPErrorStatusMetric) EvaluateMetric(ctx *gin.Context) interface{} {
    switch status := ctx.Writer.Status(); {
    case status > 399:
        return true
    default:
        return false
    }
}