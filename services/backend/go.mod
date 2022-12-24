module github.com/jmandel1027/perspex/services/backend

go 1.19

require (
	github.com/bufbuild/connect-go v1.4.1
	github.com/cockroachdb/cmux v0.0.0-20170110192607-30d10be49292
	github.com/dimiro1/health v0.0.0-20191019130555-c5cbb4d46ffc
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/jackc/pgx/v4 v4.17.2
	github.com/jmandel1027/perspex/schemas/perspex v0.0.0-00010101000000-000000000000
	github.com/jmandel1027/perspex/schemas/proto v0.0.0-00010101000000-000000000000
	github.com/prometheus/client_golang v1.14.0
	github.com/rs/cors v1.8.2
	github.com/uptrace/opentelemetry-go-extra/otelzap v0.1.17
	github.com/volatiletech/sqlboiler/v4 v4.13.0
	go.opentelemetry.io/otel v1.11.2
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.11.2
	go.opentelemetry.io/otel/sdk v1.11.2
	go.uber.org/zap v1.24.0
	google.golang.org/grpc v1.51.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bufbuild/connect-grpcreflect-go v1.0.0 // indirect
	github.com/bufbuild/connect-opentelemetry-go v0.0.0-20221216163308-175499ea7a59 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/friendsofgo/errors v0.9.2 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.14.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.13.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.12.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/opentracing/opentracing-go v1.1.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/uptrace/opentelemetry-go-extra/otelutil v0.1.17 // indirect
	github.com/volatiletech/inflect v0.0.1 // indirect
	github.com/volatiletech/strmangle v0.0.4 // indirect
	go.opentelemetry.io/otel/metric v0.34.0 // indirect
	go.opentelemetry.io/otel/trace v1.11.2 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa // indirect
	golang.org/x/net v0.2.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/genproto v0.0.0-20221202195650-67e5cbc046fd // indirect
	google.golang.org/grpc/examples v0.0.0-20221202020918-001d234e1f2d // indirect
)

replace github.com/jmandel1027/perspex/schemas/perspex => ../../schemas/perspex

replace github.com/jmandel1027/perspex/schemas/proto => ../../schemas/proto
