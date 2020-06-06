package jaeger_negroni

import (
	"os"
	"io"
	"strings"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	jaeger_config "github.com/uber/jaeger-client-go/config"
	jaeger_log "github.com/uber/jaeger-client-go/log"
	jaeger_metrics "github.com/uber/jaeger-lib/metrics"
	jaeger_client "github.com/uber/jaeger-client-go"
	opentracing "github.com/opentracing/opentracing-go"
)

var (
	enableJaegerTracing = EnableJaegerTracing()
)

// Helper function used to determine if the tracing functionality should
// be enabled based on environment variables
func EnableJaegerTracing() bool {
	
	switch enable_tracing := os.Getenv("ENABLE_JAEGER_TRACING"); strings.ToLower(enable_tracing) {
	case "false":
		return false
	case "f":
		return false
	case "0":
		return false
	default:
		return true
	}
}

type JaegerConfig struct { Host, ServiceName string; Port int }

// Helper function used to generate a Global instance of a 
// Jaeger Tracer used to report spans to the Jaeger Agent
func GetJaegerTracer(cfg JaegerConfig) io.Closer {

	log.Info(fmt.Sprintf("creating jaeger tracer for service %s at host %s at port %d", cfg.ServiceName, cfg.Host, cfg.Port))

	// crate instance of jaeger configuration
	config := jaeger_config.Configuration{
		// configure sampler
		Sampler: &jaeger_config.SamplerConfig {
			Type: jaeger_client.SamplerTypeConst,
			Param: 1,
		},
		// configurate Sampler
		Reporter: &jaeger_config.ReporterConfig{
			LogSpans: true,
			LocalAgentHostPort: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		},
	}

	// create global tracer using configuration 
	closer, err := config.InitGlobalTracer(
		cfg.ServiceName,
		jaeger_config.Logger(jaeger_log.StdLogger),
		jaeger_config.Metrics(jaeger_metrics.NullFactory),
	)

	if err != nil { log.Fatal(fmt.Sprintf("unable to create jaeger tracer: %v", err)) }

	return closer
}

// Define JaegerNegroni Middleware used to trace incoming requests using the
// Jaeger Tracing interface via the opentracing standard
func JaegerNegroni(metrics []JaegerMetric) gin.HandlerFunc {
	return func (context *gin.Context) {
		// skip middleware if tracing is disabled in environment variables
		if !enableJaegerTracing { log.Warn("jaeger tracing disabled. calls will not be traced"); context.Next(); }
		
		log.Debug(fmt.Sprintf("starting trace for route '%s'", context.FullPath()))

		// create span for each incoming request and execute in context of span
		span := opentracing.StartSpan(context.FullPath())
		defer span.Finish()

		// iterate over list of metric interfaces to generate key:value pairs in spans
		for _, metric := range(metrics) { span.SetTag(metric.MetricName(), metric.EvaluateMetric(context)) }
		
		context.Next()
	}
}
