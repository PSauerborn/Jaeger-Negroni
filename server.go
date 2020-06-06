package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	jaeger_negroni "github.com/PSauerborn/Jaeger-Negroni/jaeger_negroni"
	opentracing "github.com/opentracing/opentracing-go"
)


func test() {

	log.SetLevel(log.DebugLevel)

	// create configuration for jaeger tracer
	config := jaeger_negroni.JaegerConfig{ Host: "localhost", ServiceName: "testing-service", Port: 6831 }

	// create instance of global tracer and defer closing 
	closer := jaeger_negroni.GetJaegerTracer(config)
	defer closer.Close()
	
	for i := 0; i < 10; i++ {
		// create example span and defer closing
		span := opentracing.StartSpan("sample-span")
		span.Finish()		
	}
}

func main() {
	log.SetLevel(log.DebugLevel)

	// create new instance of router
	router := gin.New()
	router.Use(gin.Recovery())

	// create configuration for jaeger tracer
	config := jaeger_negroni.JaegerConfig{ Host: "localhost", ServiceName: "testing-service", Port: 6831 }

	// create instance of global tracer and defer closing 
	closer := jaeger_negroni.GetJaegerTracer(config)
	defer closer.Close()

	// get default metrics used to tag spans
	metrics := jaeger_negroni.DefaultJaegerMetrics()

	// apply JaegerNegroni MiddleWare for Jaeger Tracing functionality
	router.Use(jaeger_negroni.JaegerNegroni(metrics))
	router.GET("/health_check", func(context *gin.Context) { context.JSON(200, gin.H{ "http_code": 200, "message": "api running" }) })

	router.Run()
}