.PHONY: server build lint help docker clean asset test

lint:
	golint ./..

fmt:
	go fmt .

run:
	go run . server --c config/dnsx.yaml

build:
	go build -race -ldflags "-s -w -X 'dnsx/api/controller/v1.BuildTime=`date +"%Y-%m-%d %H:%M:%S"`' -X dnsx/api/controller/v1.BuildVersion=`git rev-parse --short HEAD`" -tags=jsoniter -o dnsx .

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X 'dnsx/api/controller/v1.BuildTime=`date +"%Y-%m-%d %H:%M:%S"`' -X dnsx/api/controller/v1.BuildVersion=`git rev-parse --short HEAD`" -tags=jsoniter -o dnsx.linux .
	upx dnsx.linux

docker:
	docker build -t dnsx .

clean:
	go mod tidy

test:
	go test -coverprofile c.out `go list ./... | grep -v /vendor/` -count=1 -coverpkg=`go list ./... | grep -v /vendor/`

outcov:
	go test -race -v -count=1 -coverpkg=./... -test.short -coverprofile=coverage.out -timeout=10s `go list ./... | grep -v /vendor/` -json > report.json

sonar: outcov
	sonar-scanner \
	  -Dsonar.projectKey=dnsx \
	  -Dsonar.sources=. \
	  -Dsonar.host.url=http://localhost:9000 \
	  -Dsonar.login=57104c1d43f4f9ca4b51a11c46643843cb413bc3 \
	  -Dsonar.sources.inclusions='**/*.go' \
	  -Dsonar.exclusions='doc/**,**/*_test.go,**/vendor/**,.git/**,.glide/**,asset/**,internal/asset/**' \
	  -Dsonar.tests=. -Dsonar.test.inclusions='**/*_test.go' -Dsonar.test.exclusions='**/vendor/**' \
	  -Dsonar.go.tests.reportPaths=report.json  -Dsonar.go.coverage.reportPaths=coverage.out

help:
	@echo "make: compile packages and dependencies"
	@echo "  make run: go run at server"
	@echo "  make asset: go-bindata tools"
	@echo "  make build: go build"
	@echo "  make lint: go lint ./..."
	@echo "  make clean: remove invalid go packages"