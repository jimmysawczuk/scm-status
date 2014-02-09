GOPATH := ${PWD}

define reset
	@rm -rf bin
	@mkdir -p bin
endef

define fmt
	@echo 'Running gofmt...';
	find . -type f -name "*.go" | xargs gofmt -w
endef

define build
	@echo 'Building...'

	go install scm
	go install static

	go install scm-status
endef

dev:
	@$(reset)
	@$(fmt)
	@$(build)

install:
	@sudo cp bin/scm-status /usr/local/bin/
	@sudo chmod 555 /usr/local/bin/scm-status

default: dev
