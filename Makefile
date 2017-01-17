all:
	docker run --rm -v $(PWD):/go/src/app -e GOOS=`uname -s | awk '{print tolower($$0)}'` -e GOARCH=`uname -p | sed -e 's/i386/386/'` artemave/cccv go build -v -o cccv

build-image:
	docker build -t artemave/cccv .
