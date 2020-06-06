package jaeger_negroni

import (
	"time"
	"github.com/gin-gonic/gin"
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


type RequestTimeJaegerMetric struct {}

func (metric RequestTimeJaegerMetric) EvaluateMetric(context *gin.Context) string { 
	loc, _ := time.LoadLocation("UTC")
	return time.Now().In(loc).String()
}

func (metric RequestTimeJaegerMetric) MetricName() string { return "timestamp" } 


// Helper function used to generate a map of default Jaeger Metrics
// used to set Jaeger key:value pairs within maps. The default Jaeger
// Metrics tag the spans with the URL used in the context
func DefaultJaegerMetrics() []JaegerMetric {
	return []JaegerMetric { 
		URLJaegerMetric{},
		RequestTimeJaegerMetric{},
	}
}