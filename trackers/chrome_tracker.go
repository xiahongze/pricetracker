package trackers

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

var (
	chromeTimeout = time.Second * 10
	chromeOpts    = []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.Flag("enable-automation", false),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
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
	/**
	 * Hard research to find out how to hide the automation mode
	 * https://github.com/chromedp/chromedp/issues/396
	 * https://intoli.com/blog/not-possible-to-block-chrome-headless/
	 **/
	hiddenScript = `
	Object.defineProperty(navigator, 'webdriver', {
		get: () => false,
	  });`
	hide = chromedp.ActionFunc(func(ctx context.Context) error {
		_, err := page.AddScriptToEvaluateOnNewDocument(hiddenScript).Do(ctx)
		if err != nil {
			return err
		}
		// log.Println("identifier: ", identifier.String())
		return nil
	})
)

func init() {
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
// and returns content, error
func ChromeTracker(url, xpath *string) (res string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), chromeTimeout)
	defer cancel()
	ctx, cancel = chromedp.NewExecAllocator(ctx, chromeOpts...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	log.Printf("INFO: loading %s", *url)

	err = chromedp.Run(ctx,
		hide,
		chromedp.Navigate(*url),
		chromedp.Text(*xpath, &res, chromedp.NodeVisible, chromedp.BySearch),
	)
	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\n", " ", -1)

	return
}
