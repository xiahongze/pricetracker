GO := go

build-pi:
	GOOS=linux GOARCH=arm GOARM=7 $(GO) build

build:
	$(GO) build