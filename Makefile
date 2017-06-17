PKGS := $(shell go list ./... | grep -v /vendor)

compile:
	git stash -u
	gox -output "build/{{.Dir}}_{{.OS}}_{{.Arch}}"

fmt:
	find ! -path "./vendor/*" -name "*.go" -exec gofmt -s -w {} \;

gometalinter:
	gometalinter -D gotype -D aligncheck --vendor --deadline=600s --dupl-threshold=200 -e '_string' -j 5 ./...

run-tests: setup-test-fixtures
	./test.sh

run-quick-tests: setup-test-fixtures
	go test -v $(PKGS)

test-all: gometalinter run-tests gommit doc-hunt

test-package:
	go test -race -cover -coverprofile=/tmp/strumpt github.com/antham/strumpt/$(pkg)
	go tool cover -html=/tmp/strumpt -o /tmp/strumpt.html
