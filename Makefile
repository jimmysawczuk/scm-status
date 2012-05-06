default:
	@gofmt -w src/*.go
	
	@go tool 6g -o bin/scm.6 src/scm.go src/scm_managers/*.go
	
	@go tool 6g -o bin/main.6 src/main.go
	
	@go tool 6l -o bin/main bin/main.6
	
	@rm bin/*.6
	
test:
	@@make
	@bin/main .
	
setup:
	@@make
	@bin/main setup