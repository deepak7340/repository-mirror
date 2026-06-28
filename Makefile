PACKAGE_NAME = repository-mirror
VERSION ?= 1

.PHONY: all build packages deb rpm dep clean test format vet lint

all: packages

build:
	CGO_ENABLED=0 go build -v -ldflags="-s -w" -o repository_mirror ./cmd/repo-sync

dist:
	mkdir -p dist

packages: build dist
	fpm -s dir -t deb \
		-n $(PACKAGE_NAME) \
		-v $(VERSION) \
		-a amd64 \
		--prefix /usr/local/bin \
		--description "Repository sync tool for apt/rpm mirrors" \
		--package dist/ \
		repository_mirror=/usr/local/bin/repository_mirror
	fpm -s dir -t rpm \
		-n $(PACKAGE_NAME) \
		-v $(VERSION) \
		-a x86_64 \
		--prefix /usr/local/bin \
		--description "Repository sync tool for apt/rpm mirrors" \
		--package dist/ \
		repository_mirror=/usr/local/bin/repository_mirror

deb: build dist
	fpm -s dir -t deb \
		-n $(PACKAGE_NAME) \
		-v $(VERSION) \
		-a amd64 \
		--prefix /usr/local/bin \
		--description "Repository sync tool for apt/rpm mirrors" \
		--package dist/ \
		repository_mirror=/usr/local/bin/repository_mirror

rpm: build dist
	fpm -s dir -t rpm \
		-n $(PACKAGE_NAME) \
		-v $(VERSION) \
		-a x86_64 \
		--prefix /usr/local/bin \
		--description "Repository sync tool for apt/rpm mirrors" \
		--package dist/ \
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
	rm -rf repository_mirror dist/
