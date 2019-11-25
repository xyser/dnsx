.PHONY: server build lint help docker clean asset test

run:
	go run . server --c config/config.yaml

server:
	go run . server

build:
	go build -race -ldflags "-s -w -X 'dnsx/api/controller/v1.BuildTime=`date +"%Y-%m-%d %H:%M:%S"`' -X dnsx/api/controller/v1.BuildVersion=1.0.1" -tags=jsoniter -o dnsx .

bl:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X 'dnsx/api/controller/v1.BuildTime=`date +"%Y-%m-%d %H:%M:%S"`' -X dnsx/api/controller/v1.BuildVersion=1.0.1" -tags=jsoniter -o dnsx.linux .
	upx dnsx.linux

docker:
	docker build -t dnsx .

clean:
	go mod tidy

test:
	go test $(go list ./... | grep -v /vendor/) -v -count=1 -coverpkg=./...

outcov:
	go test  -v -count=1 -coverpkg=./... -test.short -coverprofile=coverage.out -timeout=10s `go list ./... | grep -v /vendor/` -json > report.json

asset:
	go-bindata -pkg asset -o internal/asset/bindata.go asset

help:
	@echo "make: compile packages and dependencies"
	@echo "  make run: go run at server"
	@echo "  make server: go run at server"
	@echo "  make build: go build"
	@echo "  make lint: golint ./..."
	@echo "  make clean: remove object files and cached files"