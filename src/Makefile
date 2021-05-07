spec_file = server/spec/synche-server-api.yaml
build_dir = ../build

generate:
    # important: the principal should be set to schema.User
	swagger generate server \
		 --name Synche \
		 --spec $(spec_file) \
		 --principal gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema.User \
		 --flag-strategy=pflag \
		 --exclude-main \
		 -t server

	swagger generate client \
		--spec $(spec_file) \
		-c apiclient \
		-t client

build:
	go build -o $(build_dir)/synche client/main.go && \
    go build -o $(build_dir)/synche-server server/main.go
