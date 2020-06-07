package jaeger_negroni

import (
	"time"
	"strconv"
	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
)

// Define JaegerMetric interface used to create custom key:value 
// pair mappings in JaegerSpans as the API process requests
type JaegerMetric interface {
	EvaluateMetric(context *gin.Context) string
	MetricName() string
}

// Define Jaeger Metric used to tag the opentracing span with 
// the URL used in the call
type URLJaegerMetric struct {}

func (metric URLJaegerMetric) EvaluateMetric(context *gin.Context) string { return context.Request.URL.Path }

func (metric URLJaegerMetric) MetricName() string { return "url" }

// Define Jaeger Metrics used to tag the opentracing
// span with the UTC timestamp of the request time
type RequestTimeJaegerMetric struct {}

func (metric RequestTimeJaegerMetric) EvaluateMetric(context *gin.Context) string { 
	loc, _ := time.LoadLocation("UTC")
	return time.Now().In(loc).String()
}

func (metric RequestTimeJaegerMetric) MetricName() string { return "request_timestamp" } 

// Define Jaeger Metrics used to tag the opentracing
// span with the UTC timestamp of the processed time
type ProcessedTimeJaegerMetric struct {}

func (metric ProcessedTimeJaegerMetric) EvaluateMetric(context *gin.Context) string { 
	loc, _ := time.LoadLocation("UTC")
	return time.Now().In(loc).String()
}

func (metric ProcessedTimeJaegerMetric) MetricName() string { return "processed_timestamp" } 

// Define Jaeger Metrics used to tag the opentracing
// span with the request status
type RequestStatusJaegerMetric struct {}

func (metric RequestStatusJaegerMetric) EvaluateMetric(context *gin.Context) string { return strconv.Itoa(context.Writer.Status()) }

func (metric RequestStatusJaegerMetric) MetricName() string { return "request_status" } 


// Helper function used to generate a map of default Jaeger Metrics
// used to set Jaeger key:value pairs within maps. The default Jaeger
// Metrics tag the spans with the URL used in the context
func DefaultPreRequestMetrics() []JaegerMetric {
	return []JaegerMetric { 
		URLJaegerMetric{},
		RequestTimeJaegerMetric{},
	}
}

// Helper function used to generate a map of default Jaeger Metrics
// used to set Jaeger key:value pairs within maps. The default Jaeger
// Metrics tag the spans with the URL used in the context
func DefaultPostRequestMetrics() []JaegerMetric {
	return []JaegerMetric { 
		ProcessedTimeJaegerMetric{},
		RequestStatusJaegerMetric{},
	}
}

// Helper function used to set tags on opentracing spans using metrics
// passed down in array of JaegerMetric interfaces
func SetJaegerTags(context *gin.Context, metrics []JaegerMetric, span opentracing.Span) {
	
	// iterate over list of metric interfaces to generate key:value pairs in spans
	for _, metric := range(metrics) { span.SetTag(metric.MetricName(), metric.EvaluateMetric(context)) }
}