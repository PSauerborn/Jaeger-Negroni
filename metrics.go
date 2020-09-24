package jaeger_negroni

import (
    "net/http"
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
    switch ctx.Request.Method {
    case http.MethodGet:
        return "GET"
    case http.MethodPost:
        return "POST"
    case http.MethodPatch:
        return "PATCH"
    case http.MethodDelete:
        return "DELETE"
    case http.MethodPut:
        return "PUT"
    default:
        return ""
    }
}