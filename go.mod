module github.com/Bofry/lib-nsq

go 1.19

require (
	github.com/Bofry/trace v0.0.0-20230602075229-b7519a568839
	github.com/nsqio/go-nsq v1.1.0
	go.opentelemetry.io/otel v1.16.0
	go.opentelemetry.io/otel/trace v1.16.0
	golang.org/x/exp v0.0.0-20230522175609-2e198f4a06a1
)

require go.opentelemetry.io/otel/metric v1.16.0 // indirect

require (
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/joho/godotenv v1.5.1
	go.opentelemetry.io/otel/exporters/jaeger v1.16.0 // indirect
	go.opentelemetry.io/otel/sdk v1.16.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
)
