curdir=$(shell pwd)
docker_bin := $(shell command -v docker 2> /dev/null)
app := route256/loms

check: build fieldaligment-fix gofumpt lint mock test

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
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) mockgen -destination=internal/app/mocks/loms.go -package=services route256/loms/internal/app/services Loms
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) mockgen -destination=internal/app/mocks/orders_storage.go -package=services route256/loms/internal/app/services OrdersStorage
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) mockgen -destination=internal/app/mocks/stock_storage.go -package=services route256/loms/internal/app/services StocksStorage
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) mockgen -destination=internal/app/mocks/product_service.go -package=services route256/loms/internal/app/services ProductService


