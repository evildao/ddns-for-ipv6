export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE := auto

all: clean fmt package

test:
	go test -v --cover ./...

clean:
	rm -f ./bin/*

fmt:
	go fmt ./...

package:
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/ddns_darwin_amd64
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/ddns_linux_amd64
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/ddns_windows_amd64.exe
