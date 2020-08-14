APP_NAME:=gitlabtool

all: clean update test build-all

clean: 
	rm -f bin/*

update: 
	go get -u
	go mod tidy

build-all: build-linux build-arm build-mac build-win

build-linux:
	GOOS=linux go build -o bin/$(APP_NAME)

build-arm:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/$(APP_NAME)-arm

build-mac:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o bin/$(APP_NAME)-mac

build-win:
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -o bin/$(APP_NAME).exe

test: 
	go test -v ./...
