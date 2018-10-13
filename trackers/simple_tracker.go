package trackers

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
	"gopkg.in/xmlpath.v2"
)

// SimpleTracker accepts url and xpath to extract content
// and returns content/error message, ok
func SimpleTracker(url, xpath *string) (content string, ok bool) {
	defer func() {
		if !ok {
			log.Println(content)
		}
		log.Println("INFO: Found", content, "from", *url)
	}()

	xpExec, err := xmlpath.Compile(*xpath)
	if err != nil {
		content = "ERROR: failed to compile xpath %s" + *xpath
		ok = false
		return
	}

	resp, getErr := http.Get(*url)
	if getErr != nil {
		content = "ERROR: failed to fetch the website"
		ok = false
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)

	// create closure
	extractHelper := func(reader io.Reader) {
		xmlRoot, xmlErr := xmlpath.ParseHTML(reader)
		if xmlErr != nil {
			content = "ERROR: parse xml error: " + xmlErr.Error()
			ok = false
			return
		}
		content, ok = xpExec.String(xmlRoot)
		content = strings.TrimSpace(content)
		if !ok {
			content = "value not found"
			return
		}
	}

	// step 1. read directly from body
	extractHelper(bytes.NewReader(body))

	// step 2. try clean up HTML and do it again
	if !ok {
		root, err := html.Parse(bytes.NewReader(body))
		if err != nil {
			content = "ERROR: parse html" + err.Error()
			return
		}
		var b bytes.Buffer
		html.Render(&b, root)
		extractHelper(bytes.NewReader(b.Bytes()))
	}

	return
}
