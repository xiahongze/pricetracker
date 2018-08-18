package trackers

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/html"
	"gopkg.in/xmlpath.v2"
)

// SimpleTracker accepts url and xpath to extract content
// and returns content/error message, ok
func SimpleTracker(url, xpath *string) (content string, ok bool) {
	xpExec, err := xmlpath.Compile(*xpath)
	if err != nil {
		log.Printf("failed to compile xpath %s", *xpath)
		ok = false
		return
	}

	resp, getErr := http.Get(*url)
	if getErr != nil {
		log.Println("failed to fetch the website")
		ok = false
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)

	// create closure
	extractHelper := func(reader io.Reader) {
		xmlRoot, xmlErr := xmlpath.ParseHTML(reader)

		if xmlErr != nil {
			content = "parse xml error: " + xmlErr.Error()
			log.Println(content)
			ok = false
			return
		}
		if value, found := xpExec.String(xmlRoot); found {
			log.Println("Found:", value, "from", *url)
			content = value
			ok = true
		} else {
			ok = false
			content = "value not found"
		}
	}

	// step 1. read directly from body
	extractHelper(bytes.NewReader(body))

	// step 2. try clean up HTML and do it again
	if !ok {
		root, err := html.Parse(bytes.NewReader(body))
		if err != nil {
			content = "parse html error: " + err.Error()
			log.Println(content)
			return
		}
		var b bytes.Buffer
		html.Render(&b, root)
		extractHelper(bytes.NewReader(b.Bytes()))
	}

	return
}
