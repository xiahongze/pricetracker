package trackers

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
)

var (
	chromeTimeout = time.Second * 120
	chromePath    = ""
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
	ctx, cancel := context.WithTimeout(context.Background(), chromeTimeout)
	ctx1, cancel1 := context.WithTimeout(context.Background(), chromeTimeout)
	defer cancel()
	defer cancel1()
	opts := []runner.CommandLineOption{runner.Flag("headless", true)}
	if chromePath != "" {
		opts = append(opts, runner.Path(chromePath))
	}
	runnerOpt := chromedp.WithRunnerOptions(opts...)

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
	if err := c.Run(ctx1, tasks); err != nil {
		log.Println(err)
		return err.Error(), false
	}
	return strings.TrimSpace(res), true
}
