FROM golang:1.21
WORKDIR /usr/src/app
RUN apt-get -q update && apt-get -qy install netcat-traditional
RUN apt install -y protobuf-compiler
RUN go install go.uber.org/mock/mockgen@latest
RUN go install mvdan.cc/gofumpt@latest
RUN go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
RUN go install -v github.com/envoyproxy/protoc-gen-validate@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN export PATH="$PATH:$(go env GOPATH)/bin"
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /loms ./cmd/app
