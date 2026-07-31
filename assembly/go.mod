module github.com/Reensef/go-microservices-course/assembly

go 1.26.1

require (
	github.com/IBM/sarama v1.60.0
	github.com/Reensef/go-microservices-course/platform v0.0.0-00010101000000-000000000000
	github.com/Reensef/go-microservices-course/shared v0.0.0-00010101000000-000000000000
	github.com/caarlos0/env/v11 v11.4.1
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.11.1
	go.uber.org/zap v1.27.0
	golang.org/x/sync v0.22.0
	google.golang.org/protobuf v1.36.7
)

require github.com/google/go-cmp v0.7.0 // indirect

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/eapache/go-resiliency v1.7.0 // indirect
	github.com/google/uuid v1.6.0
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/klauspost/compress v1.19.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.27 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20250401214520-65e299d6c5c9 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.54.0 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/Reensef/go-microservices-course/shared => ../shared

replace github.com/Reensef/go-microservices-course/platform => ../platform
