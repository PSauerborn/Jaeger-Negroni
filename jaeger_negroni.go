package jaeger_negroni

import (
    "os"
    "io"
    "fmt"
    "strings"
    "net/http"
    "github.com/gin-gonic/gin"
    log "github.com/sirupsen/logrus"
    opentracing "github.com/opentracing/opentracing-go"
    jaeger_config "github.com/uber/jaeger-client-go/config"
    jaeger_log "github.com/uber/jaeger-client-go/log"
    jaeger_metrics "github.com/uber/jaeger-lib/metrics"
    jaeger_client "github.com/uber/jaeger-client-go"
)

var (
    enableTracing = enableJaegerTracing()
)

// function used to create new base configuration with base metrics
func Config(host, serviceName string, port int) JaegerNegroniConfig {
    jaegerConfig := JaegerConfig{JaegerHost: host, JaegerPort: port, ServiceName: serviceName}
    return JaegerNegroniConfig{
        JaegerConf: jaegerConfig,
        PreRequestMetrics: []JaegerMetric{HTTPUrlMetric{}, HTTPRequestMethodMetric{}},
        PostRequestMetrics: []JaegerMetric{},
    }
}

// function used to create new tracer instance for application
func NewTracer(cfg JaegerNegroniConfig) io.Closer {
    log.Info(fmt.Sprintf("creating new jaeger tracer for service %s for host %s:%d", cfg.JaegerConf.ServiceName, cfg.JaegerConf.JaegerHost, cfg.JaegerConf.JaegerPort))
    // create new configuration object
    config := jaeger_config.Configuration{
        Sampler: &jaeger_config.SamplerConfig{
            Type: jaeger_client.SamplerTypeConst,
            Param: 1,
        },
        Reporter: &jaeger_config.ReporterConfig{
            LogSpans: true,
            LocalAgentHostPort: fmt.Sprintf("%s:%d", cfg.JaegerConf.JaegerHost, cfg.JaegerConf.JaegerPort),
        },
    }
    // create new tracer instance using service name
    tracer, err := config.InitGlobalTracer(
        cfg.JaegerConf.ServiceName,
        jaeger_config.Logger(jaeger_log.StdLogger),
        jaeger_config.Metrics(jaeger_metrics.NullFactory),
    )

    if err != nil {
        log.Fatal(fmt.Errorf("unable to create new jaeger tracer: %v", err))
    }
    return tracer
}

// function used to determine if tracing should be enabled based on environ config
func enableJaegerTracing() bool {
    switch strings.ToLower(os.Getenv("ENABLE_JAEGER_TRACING")) {
    case "f", "false":
        return false
    default:
        return true
    }
}

// define function used to set jaeger spans in trace
func setJaegerSpans(span opentracing.Span, metrics []JaegerMetric, ctx *gin.Context) {
    for _, metric := range(metrics) {
        span.SetTag(metric.MetricName(), metric.EvaluateMetric(ctx))
    }
}

func extractSpan(req *http.Request) (*opentracing.SpanContext, error) {
    wireContext, err := opentracing.GlobalTracer().Extract(
        opentracing.HTTPHeaders,
        opentracing.HTTPHeadersCarrier(req.Header))
    if err != nil {
        return nil, err
    }
    return &wireContext, nil
}

// Jaeger Negroni tracing middleware used to trace requests
func JaegerNegroni(cfg JaegerNegroniConfig) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        if !enableTracing {
            log.Debug("jaeger tracing disabled for service. skipping trace")
            ctx.Next()
        } else {
            route := ctx.FullPath()
            log.Debug(fmt.Sprintf("starting trace for route %s", route))

            var span opentracing.Span
            parentSpan, _ := extractSpan(ctx.Request)
            if parentSpan != nil {
                log.Debug("continuing trace with parent span")
                span = opentracing.StartSpan(route, opentracing.ChildOf(*parentSpan))
            } else {
                log.Debug("starting new span")
                span = opentracing.StartSpan(route)
            }
            defer span.Finish()
            // set pre request metrics in span
            setJaegerSpans(span, cfg.PreRequestMetrics, ctx)
            ctx.Next()
            // set post request metrics in span
            setJaegerSpans(span, cfg.PostRequestMetrics, ctx)
        }
    }
}