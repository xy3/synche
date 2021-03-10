generate:
	swagger generate server \
		 --name Synche \
		 --spec api/openapi-spec/synche-server-api.yaml \
		 --principal models.Message \
		 --flag-strategy=pflag \
		 --exclude-main \
		 --with-flatten=remove-unused

build:
	go build -o ../../build


upload:
	make build
	scp -i ~/.ssh/id_rsa_old -P 6313 ../../build/server root@144.76.86.88:/root/synche/server