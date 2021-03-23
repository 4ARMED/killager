EXECUTABLE := "killager"
GITVERSION := $(shell git describe --dirty --always --tags --long)
PACKAGENAME := $(shell go list -m -f '{{.Path}}')

build: clean test
	go build -ldflags "-extldflags '-static' -X ${PACKAGENAME}/pkg/config.GitVersion=${GITVERSION}" -o ${EXECUTABLE} ./cmd/killager

build-quick: clean
	go build -ldflags "-extldflags '-static' -X ${PACKAGENAME}/pkg/config.GitVersion=${GITVERSION}" -o ${EXECUTABLE} ./cmd/killager

build-linux:
	GOOS=linux go build -ldflags "-extldflags '-static' -X ${PACKAGENAME}/pkg/config.GitVersion=${GITVERSION}" -o ${EXECUTABLE}-linux ./cmd/killager

clean:
	@rm -f ${EXECUTABLE}

test:
	go test -v ./...
