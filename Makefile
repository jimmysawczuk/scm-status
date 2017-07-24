define reset
	@rm -rf bin pkg
endef

define fmt
	@echo 'Running gofmt...';
	find . -type f -name "*.go" | xargs gofmt -w
endef

define build
	@echo 'Building...'
	go install github.com/jimmysawczuk/scm-status
endef

default: dev

dev:
	@$(reset)
	@$(fmt)
	@$(build)

production:
	@$(reset)
	@$(fmt)
	@$(build)

install:
	@sudo cp bin/scm-status /usr/local/bin/
	@sudo chmod 555 /usr/local/bin/scm-status
