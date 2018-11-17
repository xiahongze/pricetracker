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
	defer cancel()

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

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		if err := c.Shutdown(ctx); err != nil {
			log.Printf("WARN: shutdown chrome with error %s", err.Error())
		}
		cancel()

		if err := c.Wait(); err != nil {
			log.Printf("WARN: wait for chrome with error %s", err.Error())
		}
	}()

	var res string
	tasks := chromedp.Tasks{
		chromedp.Navigate(*url),
		chromedp.Text(*xpath, &res, chromedp.NodeReady, chromedp.BySearch),
	}

	// run the tasks
	if err := c.Run(ctx, tasks); err != nil {
		log.Println(err)
		return err.Error(), false
	}

	return strings.TrimSpace(res), true
}
