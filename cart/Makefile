curdir=$(shell pwd)/../
docker_bin := $(shell command -v docker 2> /dev/null)
app := route256/cart
user_id := $(shell id -u)

check: build gofumpt lint protoc-gen test

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
	$(docker_bin) run --rm -u $(user_id) -v $(curdir):/usr/src/app/ $(app) sh -c "cd cart && protoc --go_out=/usr/src/app/ --go-grpc_out=/usr/src/app/ --grpc-gateway_out=/usr/src/app/  --validate_out=lang=go:/usr/src/app/ -I /usr/src/app/cart/api/proto/google/api -I api/proto/ cart.proto"
	$(docker_bin) run --rm -u $(user_id) -v $(curdir):/usr/src/app/ $(app) sh -c "cp -r loms/api/proto/* cart/internal/app/grpc/clients/loms/pb/ && cd cart/internal/app/grpc/clients/loms/pb && protoc  --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative loms.proto"
	$(docker_bin) run --rm -u $(user_id) -v $(curdir):/usr/src/app/ $(app) sh -c "cd cart && protoc --go_out=. --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative internal/app/grpc/clients/product/pb/product.proto"

vendor-proto-google-api:
	$(docker_bin) run --rm -u $(user_id) -v $(curdir):/usr/src/app/ $(app) sh -c "\
		rm -r vendor-proto/googleapis || true && \
		rm -r cart/api/proto/google || true && \
		git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
     		https://github.com/googleapis/googleapis vendor-proto/googleapis && \
     	cd vendor-proto/googleapis && \
    	git sparse-checkout set --no-cone google/api && \
    	git checkout && \
    	cd - && \
    	mkdir -p  cart/api/proto/google && \
    	mv vendor-proto/googleapis/google/ cart/api/proto/ && \
    	rm -rf vendor-proto || true"

vendor-proto-validate:
	$(docker_bin) run --rm -u $(user_id) -v $(curdir):/usr/src/app/ $(app) sh -c "\
		rm -r vendor-proto/validate || true && \
		rm cart/api/proto/validate.proto || true && \
		git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor-proto/tmp && \
		cd vendor-proto/tmp && \
		git sparse-checkout set --no-cone validate && \
		git checkout && \
		cd - && \
		mv vendor-proto/tmp/validate/validate.proto cart/api/proto/validate.proto && \
		rm -rf vendor-proto || true"

