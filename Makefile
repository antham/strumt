PKGS := $(shell go list ./... | grep -v /vendor)

compile:
	git stash -u
	gox -output "build/{{.Dir}}_{{.OS}}_{{.Arch}}"

fmt:
	find ! -path "./vendor/*" -name "*.go" -exec gofmt -s -w {} \;

run-tests:
	./test.sh

run-quick-tests:
	go test -v $(PKGS)

test-all: run-tests

test-package:
	go test -race -cover -coverprofile=/tmp/strumt github.com/antham/strumt/$(pkg)
	go tool cover -html=/tmp/strumt -o /tmp/strumt.html
