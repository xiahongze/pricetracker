# Price-Tracker
[![Build Status](https://travis-ci.org/xiahongze/pricetracker.svg?branch=master)](https://travis-ci.org/xiahongze/pricetracker)

## Introduction

Welcome to Price-Tracker. Hope you could save money with this little app. This project is built for someone who has some technical background and wants to avoid paying extra when there is certainly a chance of getting a better deal.

Given an `xpath` and a link, one could extract price information from the website. Although `xpath` and the link could become invalid at anytime, we could still hope for the best and track down the trend for the price of interest.

## Requirements

Starting from v0.3.0, this project relies on an external API call [pushover](https://pushover.net/) which pushes notifications directly to your phones or desktops. The app is about 5USD which is worth the money in my view as you get neat messages right there on your phone. Also, you could build as many as app with their service if you like.

After setting up your pushover account, you need to create an app for this application. Then, 
you got two tokens, one for your user account and one for the app. That is it.

For the each record in datastore, you could overwrite the user and the device key such that
your messages are only delivered to specific persons or devices. By default, the app would send
messages to all devices in your account.

## Build

Dependencies:

- go 1.2+ (should work with higher versions)
- use go module

In the project directory, run

```bash
go build
```

## Setup and Serve

Once you have compiled this project, a binary executable named `pricetracker` should be sitting in your directory. There are some important environment variables that you want to setup before running it.

You need a Google Cloud account in order to save the results in `datastore` and an email account to alert you on price changes. You also need a recent `chrome` or `chromium` browser for fetching websites that are made of a single-page application.

To setup environment variables, do

```bash
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/your/gcp/credential.json"
export CHROME_TIMEOUT=60 # in seconds
$ ./pricetracker -h
Usage of ./pricetracker:
  -appToken string
        pushover app token
  -fetchLimit int
        fetch limit from google datastore (default 10)
  -port string
        server port (default "8080")
  -schdlFreq int
        schedule frequency in minutes (default 2)
  -userToken string
        pushover user token
# appToken and userToken are required!
```

## Put It in One Script

```bash
#!/bin/bash

export GOOGLE_APPLICATION_CREDENTIALS="/path/to/your/gcp/credential.json"
export CHROME_TIMEOUT=60 # in seconds

WDIR=/home/pi/apps/pricetracker

${WDIR}/pricetracker \
	-appToken APPTOKEN \
	-userToken USERTOKEN \
	-schdlFreq 10 \
	-fetchLimit 10
```

## Logrotation

We could also rotate the log by setting up a cron job like,

```
0 22 * * * /usr/sbin/logrotate /home/pi/apps/pricetracker/logrotate.conf --state /home/pi/logrotate-state
```

which can be brought up by `crontab -e`.

## RESTful APIs

### Create

POST `http://localhost:8080/create`

```json
{
	"url": "https://www.chemistwarehouse.com.au/buy/83208/braun-thermoscan-7-irt-6520",
	"xpath": "//div[@class=\"Price\"]/span",
	"expectedPrice": "$120.49",
	"name": "Braun Thermoscan 7 IRT 6520",
	"options": {
	}
}
```

### Read and Delete

POST `http://localhost:8080/read`

POST `http://localhost:8080/delete`

```json
{
	"key": "EhcKDHByaWNlLXRyYWNrcxCAgICgpPOTCa"
}
```

### Update

POST `http://localhost:8080/update`

```json
{
	"key": "EhcKDHByaWNlLXRyYWNrcxCAgICgpPOTCa",
	"options": {
	}
}
```
## License

GNU v3

Feel free to drop me a message~