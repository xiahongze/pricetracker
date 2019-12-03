package trackers

import (
	"fmt"
	"log"
	"strings"

	"github.com/antchfx/htmlquery"
)

// SimpleTracker accepts url and xpath to extract content
// and returns content, error message
func SimpleTracker(url, xpath *string) (content string, err error) {
	defer func() {
		if err == nil {
			log.Printf("INFO: Found innerText=%s", content)
		}
	}()

	log.Printf("INFO: loading %s", *url)
	doc, err := htmlquery.LoadURL(*url)
	if err != nil {
		return
	}
	elem := htmlquery.FindOne(doc, *xpath)
	if elem == nil {
		err = fmt.Errorf("WARN: failed to find element with `%s`", *xpath)
		return
	}
	content = htmlquery.InnerText(elem)
	content = strings.TrimSpace(content)
	content = strings.Replace(content, "\n", " ", -1)

	return
}
