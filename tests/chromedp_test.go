package main

import (
	"log"
	"testing"

	"github.com/xiahongze/pricetracker/trackers"
)

func TestColes(t *testing.T) {
	url := "https://shop.coles.com.au/a/a-nsw-metro-westmead/product/goldn-canola-canola-oil"
	xpath := `//span/strong[@class="product-price"]`
	price, ok := trackers.ChromeTracker(&url, &xpath)
	if !ok {
		t.Errorf("can't fetch price from %s with %s", url, xpath)
		return
	}
	log.Printf("price: %s", price)
}

func TestChemist(t *testing.T) {
	url := "https://www.chemistwarehouse.com.au/buy/1062/beconase-hayfever-nasal-spray-200-doses"
	xpath := `//div[@class="product__price"]`
	price, ok := trackers.ChromeTracker(&url, &xpath)
	if !ok {
		t.Errorf("can't fetch price from %s with %s", url, xpath)
		return
	}
	log.Printf("price: %s", price)
}
