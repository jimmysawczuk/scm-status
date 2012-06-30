default:
	@gofmt -w src/*.go src/scm_parsers/*.go

	@@mkdir -p bin
	
	@go tool 6g -o bin/scm.6 src/scm.go src/scm_parsers/*.go
	
	@go tool 6g -o bin/main.6 src/main.go
	
	@go tool 6l -o bin/scm-status bin/main.6
	
	@rm bin/*.6
	
test:
	@@make
	@bin/scm-status .
	
setup:
	@@make
	@bin/scm-status setup

install:
	@sudo cp bin/scm-status /usr/local/bin/
	@sudo chmod 555 /usr/local/bin/scm-status