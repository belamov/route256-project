curdir=$(shell pwd)
docker_bin := $(shell command -v docker 2> /dev/null)
app := route256/loms
user_id := $(shell id -u)

check: build fieldaligment-fix gofumpt lint protoc-gen test

build:
	$(docker_bin) build -t $(app) .

test:
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) go test -v -race ./...

lint:
	$(docker_bin) run --rm -v $(curdir):/app -w /app golangci/golangci-lint:latest golangci-lint run

gofumpt:
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) gofumpt -l -w .

fieldaligment-fix:
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) fieldalignment -fix ./internal/... || true

mock:
	rm -f internal/app/mocks/*
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) mockgen -destination=internal/app/services/loms_mock.go -package=services route256/loms/internal/app/services Loms
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) mockgen -destination=internal/app/services/orders_provider_mock.go -package=services route256/loms/internal/app/services OrdersProvider
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) mockgen -destination=internal/app/services/stock_provider_mock.go -package=services route256/loms/internal/app/services StocksProvider

protoc-gen:
	$(docker_bin) run --rm -u $(user_id) -v $(curdir):/usr/src/app/ $(app) protoc --go_out=/usr/src/app/ --go-grpc_out=/usr/src/app/ --grpc-gateway_out=/usr/src/app/ -I /usr/src/app/api/proto/google/api -I api/proto/ loms.proto

vendor-proto-google-api:
	$(docker_bin) run --rm -u $(user_id) -v $(curdir):/usr/src/app/ $(app) sh -c "\
		rm -r vendor-proto/googleapis || true && \
		rm -r api/proto/google || true && \
		git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
     		https://github.com/googleapis/googleapis vendor-proto/googleapis && \
     	cd vendor-proto/googleapis && \
    	git sparse-checkout set --no-cone google/api && \
    	git checkout && \
    	cd - && \
    	mkdir -p  api/proto/google && \
    	mv vendor-proto/googleapis/google/ api/proto/ && \
    	rm -rf vendor-proto || true"