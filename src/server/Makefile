generate:
    # important: the principal should be set to schema.User
	swagger generate server \
		 --name Synche \
		 --spec spec/synche-server-api.yaml \
		 --principal gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema.User \
		 --flag-strategy=pflag \
		 --exclude-main

build:
	go build -o ../../build


upload:
	make build
	scp -i ~/.ssh/id_rsa_old -P 6313 ../../build/server root@144.76.86.88:/root/synche/server