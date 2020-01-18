GOCMD=go
BINARY=awsnews
PROJECT=circa10a/go-aws-news
VERSION=0.3.0

# First target for travis ci
test:
	$(GOCMD) test -v ./...

coverage:
	$(GOCMD) test -coverprofile=c.out ./... && go tool cover -html=c.out && rm c.out

build:
	$(GOCMD) build -o $(BINARY)

docker-build:
	docker build -t $(PROJECT):$(VERSION) .

docker-run:
	docker run --rm -v $(shell pwd)/config.yaml:/config.yaml $(PROJECT):$(VERSION)

