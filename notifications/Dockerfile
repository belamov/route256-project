FROM golang:1.21
WORKDIR /usr/src/app
RUN apt-get -q update
RUN apt install -y protobuf-compiler
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
#TODO go.sum
COPY go.mod  ./
RUN go mod download && go mod verify
RUN go install go.uber.org/mock/mockgen@latest
RUN go install mvdan.cc/gofumpt@latest
RUN go install golang.org/x/tools/cmd/godoc@latest
RUN go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
RUN export PATH="$PATH:$(go env GOPATH)/bin"
COPY . .
RUN go build -v -o /notifications ./cmd/app
