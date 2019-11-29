package trackers

import (
	"fmt"
	"log"

	"github.com/antchfx/htmlquery"
)

// SimpleTracker accepts url and xpath to extract content
// and returns content/error message, ok
func SimpleTracker(url, xpath *string) (content string, ok bool) {
	defer func() {
		if !ok {
			log.Println(content)
		}
		log.Printf("INFO: Found innerText=%s", content)
	}()

	log.Printf("INFO: loading %s", *url)
	doc, err := htmlquery.LoadURL(*url)
	if err != nil {
		ok = false
		content = fmt.Sprintf("WARN: failed to load html with error %v", err)
	}
	elem := htmlquery.FindOne(doc, *xpath)
	if elem == nil {
		ok = false
		content = fmt.Sprintf("WARN: failed to find element with `%s`", *xpath)
	}
	ok = true
	content = htmlquery.InnerText(elem)

	return
}
