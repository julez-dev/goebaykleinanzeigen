package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	goebay "github.com/julez-dev/goebaykleinanzeigen"
	"golang.org/x/time/rate"
)

func main() {
	params := &goebay.SearchParam{
		Category:  goebay.Cars,
		Location:  "2750",
		Radius:    goebay.TenKM,
		PriceFrom: 1000,
		PriceTo:   6000,
		SpecificParameter: map[goebay.ParamName]string{
			goebay.CarManufacturer: "bmw",
		},
	}

	rl := rate.NewLimiter(rate.Every(60*time.Second/40), 1) // 40 reads per 60 seconds

	al := goebay.NewAdListRepo(nil)
	ar := goebay.NewAdRepo(nil)

	for {
		adList, err := al.Fetch(context.TODO(), params)

		if err != nil {
			log.Fatalln(err)
		}

		for _, ad := range adList.Items {
			if err := rl.Wait(context.TODO()); err != nil {
				log.Fatalln(err)
			}

			adItem, err := ar.Fetch(context.TODO(), ad.Link)

			if err != nil {
				log.Fatalln(err)
			}

			adItemDbg, _ := json.MarshalIndent(adItem, "", "	")
			fmt.Println(string(adItemDbg))
		}

		if adList.IsLastPage {
			break
		}

		params.Page++
	}
}
