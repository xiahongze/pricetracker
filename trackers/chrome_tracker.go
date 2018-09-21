package trackers

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
)

var (
	chromeTimeout = time.Second * 120
	chromePath    = ""
)

func init() {
	if v, ok := os.LookupEnv("CHROME_TIMEOUT"); ok {
		if vi, err := strconv.Atoi(v); err == nil {
			chromeTimeout = time.Second * time.Duration(vi)
		} else {
			log.Println("WARN: CHROME_TIMEOUT is not int but ", v)
		}
	}
	if v, ok := os.LookupEnv("CHROME_PATH"); ok {
		chromePath = v
	}
}

// ChromeTracker uses headless chrome to fetch content from given url and xpath
// and returns content/error message, ok
func ChromeTracker(url, xpath *string) (string, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), chromeTimeout)
	defer cancel()
	var runnerOpt chromedp.Option
	if chromePath == "" {
		runnerOpt = chromedp.WithRunnerOptions(
			runner.Flag("headless", true),
		)
	} else {
		runnerOpt = chromedp.WithRunnerOptions(
			runner.Path(chromePath),
			runner.Flag("headless", true),
		)
	}

	// create chrome instance
	c, err := chromedp.New(ctx, runnerOpt)
	if err != nil {
		log.Println(err)
		return err.Error(), false
	}

	var res string
	tasks := chromedp.Tasks{
		chromedp.Navigate(*url),
		chromedp.Text(*xpath, &res, chromedp.NodeVisible, chromedp.BySearch),
	}

	// run the tasks
	if err := c.Run(ctx, tasks); err != nil {
		log.Println(err)
		return err.Error(), false
	}

	return res, true
}
