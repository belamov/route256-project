curdir=$(shell pwd)
docker_bin := $(shell command -v docker 2> /dev/null)
app := route256/cart

check: build fieldaligment-fix gofumpt lint test

build:
	$(docker_bin) build -t $(app) .

test: mock
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) go test -v -race ./...

lint:
	$(docker_bin) run --rm -v $(curdir):/app -w /app golangci/golangci-lint:latest golangci-lint run

gofumpt:
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) gofumpt -l -w .

fieldaligment-fix:
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) fieldalignment -fix ./... || true

mock: build
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) mockgen -destination=internal/app/mocks/cart.go -package=services route256/cart/internal/app/services Cart
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) mockgen -destination=internal/app/mocks/cart_storage.go -package=services route256/cart/internal/app/services CartStorage
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) mockgen -destination=internal/app/mocks/product.go -package=services route256/cart/internal/app/services ProductService
	$(docker_bin) run --rm -v $(curdir):/usr/src/app/ $(app) mockgen -destination=internal/app/mocks/stock.go -package=services route256/cart/internal/app/services StocksService


