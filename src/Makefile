spec_file = synche-server-api.yaml
build_dir = ../build

generate:
    # important: the principal should be set to schema.User
	swagger generate server \
		 --name Synche \
		 --spec $(spec_file) \
		 --principal github.com/xy3/synche/src/database/schema.User \
		 --flag-strategy=pflag \
		 --exclude-main \
		 -t server

	swagger generate client \
		--spec $(spec_file) \
		-c client

build:
	go build -o $(build_dir)/synche cmd/client/main.go && \
    go build -o $(build_dir)/syncheserver cmd/server/main.go
