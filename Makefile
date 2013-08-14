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

	go install scm-status
endef

dev:
	@$(reset)
	@$(fmt)
	@$(build)

default: dev
