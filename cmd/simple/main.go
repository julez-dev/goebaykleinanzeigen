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
		Radius:    goebay.FiftyKM,
		PriceFrom: 1000,
		PriceTo:   6000,
		SpecificParameter: map[goebay.ParamName]string{
			goebay.CarManufacturer: "bmw",
		},
	}

	rl := rate.NewLimiter(rate.Every(60*time.Second/40), 1) // 40 reads per 60 seconds
	al := goebay.NewAdListRepo(nil)
	ar := goebay.NewAdRepo(nil)

	resp, err := al.Fetch(context.TODO(), params)

	if err != nil {
		log.Fatalln(err)
	}

	for _, item := range resp.Items {
		if err := rl.Wait(context.TODO()); err != nil {
			log.Fatalln(err)
		}

		// TODO ad timeouts and retries
		car, err := ar.Fetch(context.TODO(), item.Link)

		if err != nil {
			log.Fatalln(err)
		}

		carDbg, _ := json.MarshalIndent(car, "", "	")
		fmt.Println(string(carDbg))
	}
}
