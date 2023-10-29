curdir=$(shell pwd)
docker_bin := $(shell command -v docker 2> /dev/null)
docker_compose_bin := docker compose
docker_compose_yml := docker/docker-compose.yml
app := app
user_id := $(shell id -u)

check: build gofumpt lint protoc-gen sqlc test

build:
	$(docker_compose_bin) -f $(docker_compose_yml) build

up: build
	$(docker_compose_bin) -f $(docker_compose_yml) up -d

up-f: build
	$(docker_compose_bin) -f $(docker_compose_yml) up

down:
	$(docker_compose_bin) -f $(docker_compose_yml) down -v

test:
	$(docker_compose_bin) -f $(docker_compose_yml) run --rm $(app) go test -v -race ./...

lint:
	$(docker_bin) run --rm -v $(curdir):/app -w /app golangci/golangci-lint:latest golangci-lint run

gofumpt:
	$(docker_compose_bin) -f $(docker_compose_yml) run --rm $(app) gofumpt -l -w .

fieldaligment-fix:
	$(docker_compose_bin) -f $(docker_compose_yml) run --rm $(app) fieldalignment -fix ./internal/... || true

sqlc:
	$(docker_bin) run --rm -v $(curdir):/src -w /src sqlc/sqlc generate

migrate: up
	$(docker_compose_bin) -f $(docker_compose_yml) exec $(app) goose -dir internal/app/storage/migrations up
	$(docker_compose_bin) -f $(docker_compose_yml) stop

migrate-rollback: up
	$(docker_compose_bin) -f $(docker_compose_yml) exec $(app) goose -dir internal/app/storage/migrations down
	$(docker_compose_bin) -f $(docker_compose_yml) stop

mock:
	$(docker_compose_bin) -f $(docker_compose_yml) run --rm $(app) mockgen -destination=internal/app/services/loms_mock.go -package=services route256/loms/internal/app/services Loms
	$(docker_compose_bin) -f $(docker_compose_yml) run --rm $(app) mockgen -destination=internal/app/services/orders_provider_mock.go -package=services route256/loms/internal/app/services OrdersProvider
	$(docker_compose_bin) -f $(docker_compose_yml) run --rm $(app) mockgen -destination=internal/app/services/stock_provider_mock.go -package=services route256/loms/internal/app/services StocksProvider
	$(docker_compose_bin) -f $(docker_compose_yml) run --rm $(app) mockgen -destination=internal/app/services/order_event_producer_mock.go -package=services route256/loms/internal/app/services OrderEventsProducer

protoc-gen:
	$(docker_compose_bin) -f $(docker_compose_yml) run --rm $(app) protoc --go_out=. --go-grpc_out=. --grpc-gateway_out=. --validate_out=lang=go:. -I vendor.protogen -I api/proto/ loms.proto

vendor-proto-google-api:
	rm -rf vendor-proto || true
	git clone -b master --single-branch --depth=1 https://github.com/googleapis/googleapis vendor-proto/googleapis
	mkdir -p  vendor.protogen/google/api/
	mv vendor-proto/googleapis/google/api/annotations.proto vendor.protogen/google/api/annotations.proto
	mv vendor-proto/googleapis/google/api/http.proto vendor.protogen/google/api/http.proto
	rm -rf vendor-proto

vendor-proto-validate:
	rm -rf vendor-proto || true
	git clone -b main --single-branch --depth=2 --filter=tree:0 https://github.com/bufbuild/protoc-gen-validate vendor-proto/tmp
	mv vendor-proto/tmp/validate/validate.proto vendor.protogen/validate.proto
	rm -rf vendor-proto