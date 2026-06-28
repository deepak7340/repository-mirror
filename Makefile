PACKAGE_NAME = obmondo-repository-mirror
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null | sed 's/^v//')

.PHONY: all build packages deb rpm dep clean test format vet lint

all: build

build:
	CGO_ENABLED=0 go build -v -ldflags="-s -w" -o repository_mirror ./cmd/repo-sync

packages: build
	fpm -s dir -t deb \
		-n $(PACKAGE_NAME) \
		-v $(VERSION) \
		-a amd64 \
		--prefix /usr/local/bin \
		--description "Repository sync tool for apt/rpm mirrors" \
		repository_mirror=/usr/local/bin/repository_mirror
	fpm -s dir -t rpm \
		-n $(PACKAGE_NAME) \
		-v $(VERSION) \
		-a x86_64 \
		--prefix /usr/local/bin \
		--description "Repository sync tool for apt/rpm mirrors" \
		repository_mirror=/usr/local/bin/repository_mirror

deb: build
	fpm -s dir -t deb \
		-n $(PACKAGE_NAME) \
		-v $(VERSION) \
		-a amd64 \
		--prefix /usr/local/bin \
		--description "Repository sync tool for apt/rpm mirrors" \
		repository_mirror=/usr/local/bin/repository_mirror

rpm: build
	fpm -s dir -t rpm \
		-n $(PACKAGE_NAME) \
		-v $(VERSION) \
		-a x86_64 \
		--prefix /usr/local/bin \
		--description "Repository sync tool for apt/rpm mirrors" \
		repository_mirror=/usr/local/bin/repository_mirror

lint:
	@golangci-lint run --issues-exit-code=1

format:
	@go fmt ./...

vet:
	@go vet ./...

test:
	@go test -v ./...

dep:
	@go get -v ./...

clean:
	@go clean
	rm -f repository_mirror *.deb *.rpm
