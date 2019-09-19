VERSION=$(shell echo "$$(git rev-parse --abbrev-ref HEAD)-$$(git rev-parse --short HEAD)")
GOBUILD=go build -ldflags "-X github.com/GivenZeng/adn.version=$(VERSION)"

version:
	@echo ${VERSION}

pull:
	git pull

clean:
	@find bin/ -type f -executable -exec rm -v {} \;


goip:
	$(GOBUILD) -o bin/$@ github.com/GivenZeng/goip

run:
	cd bin && ./goip

goip_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make goip