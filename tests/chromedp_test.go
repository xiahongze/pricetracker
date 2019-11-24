package main

import (
	"log"
	"testing"

	"github.com/xiahongze/pricetracker/trackers"
)

func TestChromedp(t *testing.T) {
	url := "https://shop.coles.com.au/a/a-nsw-metro-westmead/product/goldn-canola-canola-oil"
	xpath := `//span/strong[@class="product-price"]`
	price, ok := trackers.ChromeTracker(&url, &xpath)
	if !ok {
		t.Errorf("can't fetch price from %s with %s", url, xpath)
		return
	}
	log.Printf("price: %s", price)
}
