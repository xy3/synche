spec_file = synche-server-api2.yaml
build_dir = ../build
schema_dir = schema

generate:
	swagger generate spec -m -w $(schema_dir) -i $(spec_file) -o $(spec_file)

	swagger generate server \
		 --name Synche \
		 --spec $(spec_file) \
		 --principal github.com/xy3/synche/src/server/schema.User \
		 --flag-strategy=pflag \
		 --exclude-main \
		 --target server

	swagger generate client \
		--spec $(spec_file) \
		-c client


build:
	go build -o $(build_dir)/synche cmd/client/main.go && \
    go build -o $(build_dir)/synched cmd/server/main.go
