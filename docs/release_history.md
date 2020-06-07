## Jaeger Negroni Release Notes

The following release notes are intended to provide a brief overview of the significant
changes made across different versions of the Jaeger Negroni project

### __Jaeger Negroni Version__ - 0.0.2a
***
__Date Updated__: June 8th 2020 <br>
__Developer(s)__: Pascal Sauerborn <br>
__Work Item Type__: Development

#### __Release Notes__

Added `pre-request` and `post-request` metrics to the `JaegerNegroni` Middleware

#### Files Modified

server/udp_server/
* modified `jaeger_metrics.go` to incorporate `ProcessedTimeJaegerMetric` and `RequestStatusJaegerMetric`
* modified `JaegerNegroni` middleware to accept pre and post request Jaeger Metrics

### __Jaeger Negroni Version__ - 0.0.1a
***
__Date Updated__: June 7th 2020 <br>
__Developer(s)__: Pascal Sauerborn <br>
__Work Item Type__: Development

#### __Release Notes__

Added initial code base containing `gin-gonic` middleware and default jaeger metrics

#### Files Added

`/`
* added `server.go` containing sample application

`/jaeger_negroni`
* added `helpers.go` file containing helper functions used by jaeger middleware
* added `jaeger_metrics.go` file containing `JaegerMetric` interface as well as interfaces for `URLJaegerMetrics` and `RequestTimeJaegerMetrics`
* added `jaeger_negroni.go` file containing `JaegerNegroni` Middleware


