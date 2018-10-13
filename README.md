# Price-Tracker

## Introduction

Welcome to Price-Tracker. Hope you could save money with this little app. This project is built for someone who has some technical background and wants to avoid paying extra when there is certainly a chance of getting a better deal.

Given an `xpath` and a link, one could extract price information from the website. Although `xpath` and the link could become invalid at anytime, we could still hope for the best and track down the trend for the price of interest.

## Build

Dependencies:

- go 1.8.3 (should work with higher versions)
- dep v0.5.0
- and those listed in `Gopkg.toml`

In the project directory, run

```bash
make install-deps build
```

## Setup and Serve

Once you have compiled this project, a binary executable named `pricetracker` should be sitting in your directory. There are some important environment variables that you want to setup before running it.

You need a Google Cloud account in order to save the results in `datastore` and an email account to alert you on price changes. You also need a recent `chrome` or `chromium` browser for fetching websites that are made of a single-page application.

To setup environment variables, do

```bash
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/your/gcp/credential.json"
export EMAIL_USERNAME="user"
export EMAIL_PASSWORD="password"
export EMAIL_HOST="smtp.gmail.com"
export EMAIL_FROM="user@gmail.com"
export SCHEDULE_FREQ=10 # in minutes
export CHROME_PATH="/usr/bin/chromium-browser"
export CHROME_TIMEOUT=60 # in seconds
```

and then simply run, `./pricetracker`. Well done! You make it.

## Daemonization

You can do better by making this little app run as a daemon in Linux by,

```bash
#!/bin/bash

WDIR=/home/pi/apps/pricetracker

source ${WDIR}/.env && daemon --name="pricetracker" --output=${WDIR}/log ${WDIR}/pricetracker
```

where `.env` is the environment variables (saved as a file) we set in the last step.

Of course, you need `daemon` for this script to work.

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
		"email": "user@gmail.com"
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
		"email": "user@gmail.com"
	}
}
```
## License

GNU v3

Feel free to drop me a message~