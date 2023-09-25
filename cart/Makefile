curdir=$(shell pwd)
docker_bin := $(shell command -v docker 2> /dev/null)
app := route256/cart

check: build fieldaligment-fix gofumpt lint test

build:
	$(docker_bin) build -t $(app) .

test:
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) go test -v -race ./...

lint:
	$(docker_bin) run --rm -v $(curdir):/app -w /app golangci/golangci-lint:latest golangci-lint run

gofumpt:
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) gofumpt -l -w .

fieldaligment-fix:
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) fieldalignment -fix ./... || true

mock:
	rm -f internal/app/mocks/*
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) mockgen -destination=internal/app/services/cart_mock.go -package=services route256/cart/internal/app/services Cart
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) mockgen -destination=internal/app/services/cart_provider_mock.go -package=services route256/cart/internal/app/services CartProvider
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) mockgen -destination=internal/app/services/product_mock.go -package=services route256/cart/internal/app/services ProductService
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) mockgen -destination=internal/app/services/loms_mock.go -package=services route256/cart/internal/app/services LomsService


