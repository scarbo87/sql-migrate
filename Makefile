.PHONY: all build build_osx

build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sql-migrate-cobra_linux_amd64 main.go

build_osx:
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o sql-migrate-cobra_darwin_amd64 main.go