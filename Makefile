SERVICE = atomica-blog
PKG_LIST = $(shell go list ./... | grep -v mock)
COMMIT := $(shell sh -c 'git rev-parse --short HEAD')
clean:
	rm -rf ./bin

test:
	go test -v $(PKG_LIST)

lint:
	golint -set_exit_status $(PKG_LIST)

build-osx: clean generate-swagger
    GOOS=darwin GOARCH=amd64 go build -v -a -tags aws_lambda -o bin/$(SERVICE)-api -a --ldflags "-w \
    -X github.com/asaberwd/atomica-blog/build.GitCommit=$(COMMIT)" cmd/api/main.go

build: clean generate-swagger
	GOOS=linux GOARCH=amd64 go build -v -a -tags aws_lambda -o bin/$(SERVICE)-api -a --ldflags "-w \
	-X github.com/asaberwd/atomica-blog/build.GitCommit=$(COMMIT)" cmd/api/main.go

validate-swagger:
	swagger validate ./swagger/api.yaml

generate-swagger: validate-swagger
	swagger generate server \
		--target=./swagger \
		--spec=./swagger/api.yaml -A atomica-blog-service

local-run:
	go run ./swagger/cmd/atomica-blog-service-server/main.go --scheme http --port=8050
