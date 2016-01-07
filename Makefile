define reset
	@rm -rf bin pkg
endef

define fmt
	@echo 'Running gofmt...';
	find . -type f -name "*.go" | xargs gofmt -w
endef

define build
	@echo 'Building...'
	go install scm-status/cmd/scm-status
endef

define snapshot
	@echo 'Snapshotting...'
	@go get github.com/jimmysawczuk/go-binary
	@git tag --contains HEAD | go-binary -f="getVersion" -p="scm" -out="scm/version.go"
endef

default: dev

snapshot:
	@$(snapshot)

dev:
	@$(reset)
	@$(fmt)
	@$(build)

production:
	@$(reset)
	@$(fmt)
	@$(snapshot)
	@$(build)

install:
	@sudo cp bin/scm-status /usr/local/bin/
	@sudo chmod 555 /usr/local/bin/scm-status
