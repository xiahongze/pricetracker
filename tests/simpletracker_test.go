package main

import (
	"log"
	"testing"

	"github.com/xiahongze/pricetracker/trackers"
)

func TestChemistSimple(t *testing.T) {
	url := "https://www.chemistwarehouse.com.au/buy/1062/beconase-hayfever-nasal-spray-200-doses"
	xpath := `//span[@class="product__price"] | //div[@class="product__price"]`
	price, err := trackers.SimpleTracker(&url, &xpath)
	if err != nil {
		t.Errorf("can't fetch price from %s with %s error: %v", url, xpath, err)
		return
	}
	log.Printf("price: %s", price)
}
