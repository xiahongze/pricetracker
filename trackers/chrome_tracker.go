package trackers

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

var (
	chromeTimeout = time.Second * 120
	chromePath    = ""
	chromeOpts    = []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-breakpad", true),
		chromedp.Flag("disable-client-side-phishing-detection", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-features", "site-per-process,TranslateUI,BlinkGenPropertyTrees"),
		chromedp.Flag("disable-hang-monitor", true),
		chromedp.Flag("disable-ipc-flooding-protection", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-prompt-on-repost", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-sync", true),
	}
)

func init() {
	if v, ok := os.LookupEnv("CHROME_PATH"); ok {
		chromePath = v
	}
	if v, ok := os.LookupEnv("CHROME_TIMEOUT"); ok {
		vi, err := strconv.Atoi(v)
		if err != nil {
			log.Println("WARN: CHROME_TIMEOUT is not int but ", v)
			return
		}
		chromeTimeout = time.Second * time.Duration(vi)
	}
}

// ChromeTracker uses headless chrome to fetch content from given url and xpath
// and returns content/error message, ok
func ChromeTracker(url, xpath *string) (string, bool) {
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), chromeOpts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var res string
	err := chromedp.Run(ctx, chromedp.Navigate(*url), chromedp.Text(*xpath, &res, chromedp.NodeVisible, chromedp.BySearch))

	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimSpace(res), true
}
