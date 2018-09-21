GO := go

install-dependencies:
	$(GO) get -u gopkg.in/xmlpath.v2
	$(GO) get -u golang.org/x/net/html
	$(GO) get -u golang.org/x/sys/unix
	$(GO) get -u github.com/chromedp/chromedp
	$(GO) get -u cloud.google.com/go/datastore
	$(GO) get -u github.com/labstack/echo/...

build-pi:
	GOOS=linux GOARCH=arm GOARM=7 $(GO) build