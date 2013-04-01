OS = $(GOOS)
ARCH = $(GOARCH)

OUTFILE = bin/scm-status

default: fmt build link

fmt:
	@echo "Running gofmt..."
	@gofmt -w src/*.go src/scm_parsers/*.go

build:
	@mkdir -p bin

	@echo "Compiling..."
	@GOOS=$(OS) GOARCH=$(ARCH) go tool 6g -o bin/scm.6 src/scm.go src/scm_parsers/*.go
	@GOOS=$(OS) GOARCH=$(ARCH) go tool 6g -o bin/main.6 src/main.go

link: fmt build
	@echo "Linking..."
	@GOOS=$(OS) GOARCH=$(ARCH) go tool 6l -o $(OUTFILE) bin/main.6
	@rm bin/*.6

test: default
	@bin/scm-status .

dist:
	@mkdir -p dist

	@@make -s OS=linux ARCH=amd64 OUTFILE=dist/scm-status.linux-amd64
	@@make -s OS=linux ARCH=386 OUTFILE=dist/scm-status.linux-386

	@@make -s OS=darwin ARCH=amd64 OUTFILE=dist/scm-status.darwin-amd64

	@@make -s OS=windows ARCH=amd64 OUTFILE=dist/scm-status.windows-amd64.exe
	@@make -s OS=windows ARCH=386 OUTFILE=dist/scm-status.windows-386.exe

install: default
	@sudo cp bin/scm-status /usr/local/bin/
	@sudo chmod 555 /usr/local/bin/scm-status
