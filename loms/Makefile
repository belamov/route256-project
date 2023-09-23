curdir=$(shell pwd)
docker_bin := $(shell command -v docker 2> /dev/null)
app := route256/loms

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
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) mockgen -destination=internal/app/domain/services/loms_mock.go -package=services route256/loms/internal/app/domain/services Loms
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) mockgen -destination=internal/app/domain/services/orders_provider_mock.go -package=services route256/loms/internal/app/domain/services OrdersProvider
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) mockgen -destination=internal/app/domain/services/stock_provider_mock.go -package=services route256/loms/internal/app/domain/services StocksProvider


