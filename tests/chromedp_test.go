package main

import (
	"log"
	"testing"

	"github.com/xiahongze/pricetracker/trackers"
)

func TestColes(t *testing.T) {
	// url := "https://shop.coles.com.au/a/national/product/goldn-canola-canola-oil"
	// url := "https://shop.coles.com.au/a/national/product/grinders-organic-coffee-beans"
	url := "https://shop.coles.com.au/a/national/product/lavazza-ground-coffee-torino-qualita-oro-4965938p"
	xpath := `//span/strong[@class="product-price"] | //*[@id="main-content-inside"]/div[2]/div/header/div[3]/div/span[1]`
	price, err := trackers.ChromeTracker(&url, &xpath)
	if err != nil {
		t.Errorf("can't fetch price from %s with %s error: %v", url, xpath, err)
		return
	}
	log.Printf("price: %s", price)
}

func TestChemist(t *testing.T) {
	url := "https://www.chemistwarehouse.com.au/buy/1062/beconase-hayfever-nasal-spray-200-doses"
	xpath := `//span[@class="product__price"] | //div[@class="product__price"]`
	price, err := trackers.ChromeTracker(&url, &xpath)
	if err != nil {
		t.Errorf("can't fetch price from %s with %s error: %v", url, xpath, err)
		return
	}
	log.Printf("price: %s", price)
}
