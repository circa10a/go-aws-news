GOCMD=go
BINARY=awsnews
BUILD_FLAGS=-ldflags="-s -w"
PROJECT=circa10a/go-aws-news
VERSION=1.2.0

test:
	$(GOCMD) test -v ./... -coverprofile=coverage.txt
	@if ! command -v golangci-lint 1>/dev/null; then\
		echo "Need to install golangci-lint";\
		exit 1;\
	fi;\
	golangci-lint run

coverage:
	$(GOCMD) test -coverprofile=coverage.txt ./... && go tool cover -html=coverage.txt

build:
	$(GOCMD) build -o $(BINARY)

compile:
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o bin/linux/amd64/$(BINARY)
	GOOS=linux GOARCH=arm go build $(BUILD_FLAGS) -o bin/linux/arm/$(BINARY)
	GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o bin/linux/arm64/$(BINARY)
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o bin/darwin/amd64/$(BINARY)
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o bin/windows/amd64/$(BINARY).exe

package:
	GOOS=linux GOARCH=amd64 tar -czf bin/$(BINARY)_$(VERSION)_linux_amd64.tar.gz -C bin/linux/amd64 $(BINARY)
	GOOS=linux GOARCH=arm tar -czf bin/$(BINARY)_$(VERSION)_linux_arm.tar.gz -C bin/linux/arm $(BINARY)
	GOOS=linux GOARCH=arm64 tar -czf bin/$(BINARY)_$(VERSION)_linux_arm64.tar.gz -C bin/linux/arm64 $(BINARY)
	GOOS=darwin GOARCH=amd64 tar -czf bin/$(BINARY)_$(VERSION)_darwin_amd64.tar.gz -C bin/darwin/amd64 $(BINARY)
	GOOS=windows GOARCH=amd64 zip -j bin/$(BINARY)_$(VERSION)_windows_amd64.zip bin/windows/amd64/$(BINARY).exe

clean:
	$(GOCMD) clean
	rm -rf $(BINARY) bin/

release: clean compile package

docker-build:
	docker build -t $(PROJECT):$(VERSION) .

docker-run:
	docker run --rm -v $(shell pwd)/config.yaml:/config.yaml $(PROJECT):$(VERSION)

docker-release: docker-build
docker-release:
	echo "${DOCKER_PASSWORD}" | docker login -u ${DOCKER_USERNAME} --password-stdin
	docker push $(PROJECT):$(VERSION)

lambda-build:
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o bin/linux/amd64/$(BINARY)

lambda-package: lambda-build
	zip -j bin/lambda.zip bin/linux/amd64/$(BINARY)

lambda-run: lambda-build
	docker run --rm -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY -e AWS_REGION=us-east-1 \
		-v ${PWD}/bin/linux/amd64:/var/task:ro,delegated lambci/lambda:go1.x $(BINARY)
