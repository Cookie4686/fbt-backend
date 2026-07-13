module fbt/backend

go 1.25.7

require (
	github.com/caarlos0/env/v11 v11.3.1
	github.com/jackc/pgx/v5 v5.10.0
	github.com/joho/godotenv v1.5.1
	github.com/pquerna/otp v1.5.0
	github.com/stretchr/testify v1.11.1
	golang.org/x/crypto v0.52.0
)

require (
	buf.build/go/protovalidate v1.2.0 // indirect
	cel.dev/expr v0.25.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/google/cel-go v0.28.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20260410095643-746e56fc9e2f // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260226221140-a57be14db171 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260420184626-e10c466a9529 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.11-20260415201107-50325440f8f2.1
	connectrpc.com/connect v1.20.0
	connectrpc.com/validate v0.6.0
	github.com/pressly/goose/v3 v3.27.2
	github.com/wneessen/go-mail v0.7.3
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.81.1
	google.golang.org/protobuf v1.36.11
)

tool (
	connectrpc.com/connect/cmd/protoc-gen-connect-go
	google.golang.org/protobuf/cmd/protoc-gen-go
)
