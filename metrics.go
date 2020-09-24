package jaeger_negroni

import (
    "github.com/gin-gonic/gin"
)


type JaegerMetric interface {
    MetricName() string
    EvaluateMetric(ctx *gin.Context) string
}

type HTTPUrlMetric struct{}

func(metric HTTPUrlMetric) MetricName() string {
    return "http.url"
}

func(metric HTTPUrlMetric) EvaluateMetric(ctx *gin.Context) string {
    return ctx.FullPath()
}

type HTTPRequestMethodMetric struct{}

func(metric HTTPRequestMethodMetric) MetricName() string {
    return "http.method"
}

func(metric HTTPRequestMethodMetric) EvaluateMetric(ctx *gin.Context) string {
    return ctx.Request.Method
}