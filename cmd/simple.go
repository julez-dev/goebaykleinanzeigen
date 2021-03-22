package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/julez-dev/goebay"
	"golang.org/x/time/rate"
)

func main() {
	params := &goebay.SearchParam{
		Category:  goebay.Cars,
		Location:  "3331",
		Radius:    goebay.FiftyKM,
		PriceFrom: 1000,
		PriceTo:   4000,
		SpecificParameter: map[goebay.ParamName]string{
			goebay.CarManufacturer: "bmw",
		},
	}

	rl := rate.NewLimiter(rate.Every(60*time.Second/40), 1) // 40 reads per 60 seconds

	al := goebay.NewAdListRepo(nil)
	ar := goebay.NewAdRepo(nil)

	resp, err := al.Fetch(context.Background(), params)

	if err != nil {
		log.Fatalln(err)
	}

	for _, item := range resp.Items {
		if err := rl.Wait(context.Background()); err != nil {
			log.Fatalln(err)
		}

		car, err := ar.Fetch(context.Background(), item.Link)

		if err != nil {
			log.Fatalln(err)
		}

		carDbg, _ := json.MarshalIndent(car, "", "	")
		fmt.Println(string(carDbg))
	}
}
