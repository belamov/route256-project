curdir=$(shell pwd)
docker_bin := $(shell command -v docker 2> /dev/null)
docker_compose_bin := docker compose
docker_compose_yml := docker/docker-compose.yml
app := app
user_id := $(shell id -u)

check: build gofumpt lint protoc-gen test

build:
	$(docker_compose_bin) -f $(docker_compose_yml) build

up: build
	$(docker_compose_bin) -f $(docker_compose_yml) up -d

down:
	$(docker_compose_bin) -f $(docker_compose_yml) down -v

stop:
	$(docker_compose_bin) -f $(docker_compose_yml) stop

test:
	$(docker_compose_bin) -f $(docker_compose_yml) run --rm $(app) go test -v -race ./...

lint:
	$(docker_bin) run --rm -v $(curdir):/app -w /app golangci/golangci-lint:latest golangci-lint run

gofumpt:
	$(docker_compose_bin) -f $(docker_compose_yml) run --rm $(app) gofumpt -l -w .

fieldaligment-fix:
	$(docker_compose_bin) -f $(docker_compose_yml) run --rm $(app) fieldalignment -fix ./internal/... || true

sqlc:
	$(docker_bin) run --rm -u $(user_id) -v $(curdir):/src -w /src sqlc/sqlc generate

create-migration:
	$(docker_compose_bin) -f $(docker_compose_yml) run --rm -u $(user_id) $(app) goose -dir internal/app/storage/migrations create $(name) sql

migrate: up
	$(docker_compose_bin) -f $(docker_compose_yml) exec $(app) goose -dir internal/app/storage/migrations up
	$(docker_compose_bin) -f $(docker_compose_yml) stop

mock:
	$(docker_compose_bin) -f $(docker_compose_yml) run --rm $(app) mockgen -destination=internal/app/services/cart_mock.go -package=services route256/cart/internal/app/services Cart
	$(docker_compose_bin) -f $(docker_compose_yml) run --rm $(app) mockgen -destination=internal/app/services/cart_provider_mock.go -package=services route256/cart/internal/app/services CartProvider
	$(docker_compose_bin) -f $(docker_compose_yml) run --rm $(app) mockgen -destination=internal/app/services/product_mock.go -package=services route256/cart/internal/app/services ProductService
	$(docker_compose_bin) -f $(docker_compose_yml) run --rm $(app) mockgen -destination=internal/app/services/loms_mock.go -package=services route256/cart/internal/app/services LomsService

protoc-gen:
	# pb for cart
	$(docker_compose_bin) -f $(docker_compose_yml) run --rm $(app) protoc --go_out=. --go-grpc_out=. --grpc-gateway_out=. --validate_out=lang=go:. -I api/proto -I vendor.protogen cart.proto
	# pb for loms client with loms proto
	mkdir -p internal/app/grpc/clients/loms/pb
	$(docker_compose_bin) -f $(docker_compose_yml) run --rm $(app) protoc  --go_out=internal/app/grpc/clients/loms/pb/ --go-grpc_out=internal/app/grpc/clients/loms/pb/ --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative -I internal/app/grpc/clients/loms/pb/ -I vendor.protogen internal/app/grpc/clients/loms/pb/loms.proto
	# pb for product
	$(docker_compose_bin) -f $(docker_compose_yml) run --rm $(app) protoc --go_out=. --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative internal/app/grpc/clients/product/pb/product.proto

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