# goebaykleinanzeigen

## _Disclaimer_

This project is in early development and **will** be buggy.

I am not associated with eBay Kleinanzeigen or eBay in any way.

Use this code at your own risk.

## Introduction

goebaykleinanzeigen is an experimental library for scraping eBay Kleinanzeigens data written in go.

eBay Kleinanzeigen should not be confused with the main site eBay.

There is a runnable basic example for scraping cars in the cmd directory.

This library is written with the intend to scrape car listings but all categories should work and will be added in the future.

## Why webscraping?

eBay Kleinanzeigen offers a private API which is only available for their official partners.
If you need to use this library you are probably **not** a official partner.

Sadly this makes accessing eBay Kleinanzeigens data a ~~little~~ lot harder. Especially since we are limited to about 40 req/minute.
And of course parsing HTML is pretty error prone.

## eBay Kleinanzeigen URLs

In the following we will discuss the relevant URLs for our little scraper

### The search URL

First we need to list all the ads.

The basic URL for this may look like this: <https://www.ebay-kleinanzeigen.de/l3331r50>. Here we can split up the parameter in 2 parts: l3331 and r50

- Where the 3331 in l3331 is some kind of internal id? for Berlin.
- Where the 50 in r50 is the distance from the center in kilometers.

This is pretty boring and generic so lets specify a specific category in our URL: <https://www.ebay-kleinanzeigen.de/l3331r50c216>

- We added the parameter c216 to our URL! In this case c is probably a prefix for categoryID and 216 is the ID for cars!

So now we only search for cars but this is still a little bit to generic if you ask me. We don't want to spend too much money on our new car so let's add a specific
price range to our query: <https://www.ebay-kleinanzeigen.de/preis:1000:4000/l3331r50c216>

- Now we only search for cars in the specific price range from 1000 euro to 4000 euro (preis is the german word for price).

But we don't want a Volkswagen, we only want to search for BMWs. Luckily we can tell eBay to only look for BMWs! Our URL now looks like this: <https://www.ebay-kleinanzeigen.de/preis:1000:4000/l3331r50c216+autos.marke_s:bmw>.

- Options specific for a category are appended in the following format: +option_name:option_value

You don't really need to worry about this since the `SearchParam´ struct abstracts the query building away. The SearchParam for this example may look like this:

```go
param := &goebay.SearchParam{
		Category:  goebay.Cars,
		CountryID: "3331",
		Radius:    goebay.FiftyKM,
		PriceFrom: 1000,
		PriceTo:   4000,
		SpecificParameter: map[goebay.ParamName]string{
			goebay.CarManufacturer: "bmw",
		},
	}
```

More specific parameters will be added to `specificparams.go` in the future.

### The ad URL

This is pretty straight forward. The URL for an ad looks like this: <https://www.ebay-kleinanzeigen.de/s-anzeige/{ADID}>

### LocationID for a Location

The internal ID for a location can be fetched by this URL (Berlin in this case): <https://www.ebay-kleinanzeigen.de/s-ort-empfehlungen.json?query=Berlin>

## Ratelimit

Currently the max seems to be about 40 req/minute. This is currently not enforced by the library

## Localisation

Keep in mind that eBay Kleinanzeigen is a German Site so everything will be in German.

Prices are integers in euro there are no cents on eBay Kleinanzeigen(? I did not confirm this).
