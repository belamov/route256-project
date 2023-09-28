curdir=$(shell pwd)/../
docker_bin := $(shell command -v docker 2> /dev/null)
app := route256/cart
user_id := $(shell id -u)

check: build fieldaligment-fix gofumpt lint protoc-gen test

build:
	$(docker_bin) build -t $(app) -f Dockerfile ../

test:
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) sh -c "cd cart && go test -v -race ./..."

lint:
	$(docker_bin) run --rm -v $(curdir):/app -w /app/cart golangci/golangci-lint:latest golangci-lint run

gofumpt:
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) sh -c "cd cart && gofumpt -l -w ."

fieldaligment-fix:
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) sh -c "cd cart && fieldalignment -fix ./internal/... || true"

mock:
	rm -f internal/app/mocks/*
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) sh -c "cd cart && mockgen -destination=internal/app/services/cart_mock.go -package=services route256/cart/internal/app/services Cart"
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) sh -c "cd cart && mockgen -destination=internal/app/services/cart_provider_mock.go -package=services route256/cart/internal/app/services CartProvider"
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) sh -c "cd cart && mockgen -destination=internal/app/services/product_mock.go -package=services route256/cart/internal/app/services ProductService"
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) sh -c "cd cart && mockgen -destination=internal/app/services/loms_mock.go -package=services route256/cart/internal/app/services LomsService"

protoc-gen:
	$(docker_bin) run --rm -u $(user_id) -v $(curdir):/usr/src/app/ $(app) sh -c "cd cart && protoc --go_out=. --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative api/proto/cart.proto"

