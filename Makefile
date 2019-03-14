USER=$(shell whoami)
HEAD=$(shell ([ -n "$${CI_TAG}" ] && echo "$$CI_TAG" || exit 1) || git describe --tags 2> /dev/null || git rev-parse --short HEAD)
STAMP=$(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
DIRTY=$(shell test $(shell git status --porcelain | wc -l) -eq 0 || echo '(dirty)')


LDFLAGS="-X main.buildStamp=$(STAMP) -X main.buildUser=$(USER) -X main.buildHash=$(HEAD) -X main.buildDirty=$(DIRTY)"
all: install

build: darwin64 linux64

test:
	go test './...'

clean:
	-rm -f chksutil
	-rm -rf release
	go clean -i ./cmd/chksutil

install:
	go install -ldflags $(LDFLAGS) ./cmd/chksutil

darwin64:
	env GOOS=darwin GOARCH=amd64 go clean -i  ./cmd/chksutil
	env GOOS=darwin GOARCH=amd64 go build -ldflags $(LDFLAGS) -o release/darwin64/chksutil ./cmd/chksutil

linux64:
	env GOOS=linux GOARCH=amd64 go clean -i
	env GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -o release/linux64/chksutil ./cmd/chksutil

.PHONY: release
release: clean build
	mkdir release/dist
	gzip < release/darwin64/chksutil > 'release/dist/chksutil.darwin_amd64.$(HEAD)$(DIRTY).gz' 
	gzip <  release/linux64/chksutil > 'release/dist/chksutil.linux_amd64.$(HEAD)$(DIRTY).gz'
