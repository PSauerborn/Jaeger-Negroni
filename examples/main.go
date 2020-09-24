package main

import (
    "github.com/gin-gonic/gin"
    log "github.com/sirupsen/logrus"
    jaeger_negroni "github.com/PSauerborn/jaeger-negroni"
)

func main() {
    log.SetLevel(log.DebugLevel)

    // create new instance of router
    router := gin.New()
    router.Use(gin.Recovery())
    // create new config object
    config := jaeger_negroni.NewConfig("localhost", "test-service", 6831)
    tracer := jaeger_negroni.NewTracer(config)
    defer tracer.Close()

    // apply JaegerNegroni MiddleWare for Jaeger Tracing functionality
    router.Use(jaeger_negroni.JaegerNegroni(config))
    router.GET("/health_check", func(context *gin.Context) { context.JSON(200, gin.H{ "http_code": 200, "message": "api running" }) })

    router.Run(":10999")
}