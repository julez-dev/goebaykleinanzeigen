package goebaykleinanzeigen

import (
	"testing"
)

func Test_fmtCategory(t *testing.T) {

	tests := []struct {
		name string
		sp   *SearchParam
		want string
	}{
		{
			name: "empty-category",
			sp:   &SearchParam{},
			want: "",
		},
		{
			name: "category-car",
			sp:   &SearchParam{Category: Cars},
			want: "/c" + string(Cars),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.sp.fmtCategory(); got != tt.want {
				t.Errorf("SearchParam.fmtCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fmtPrice(t *testing.T) {

	tests := []struct {
		name string
		sp   *SearchParam
		want string
	}{
		{
			name: "empty-price",
			sp:   &SearchParam{},
			want: "",
		},
		{
			name: "price-only-from",
			sp:   &SearchParam{PriceFrom: 1000},
			want: "/preis:1000:",
		},
		{
			name: "price-only-to",
			sp:   &SearchParam{PriceTo: 1000},
			want: "/preis::1000",
		},
		{
			name: "price-both",
			sp:   &SearchParam{PriceFrom: 1000, PriceTo: 6000},
			want: "/preis:1000:6000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.sp.fmtPrice(); got != tt.want {
				t.Errorf("SearchParam.fmtPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fmtPage(t *testing.T) {

	tests := []struct {
		name string
		sp   *SearchParam
		want string
	}{
		{
			name: "empty-page",
			sp:   &SearchParam{},
			want: "",
		},
		{
			name: "page-1",
			sp:   &SearchParam{Page: 1},
			want: "/seite:1",
		},
		{
			name: "page-mius-1",
			sp:   &SearchParam{Page: -1},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.sp.fmtPage(); got != tt.want {
				t.Errorf("SearchParam.fmtPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fmtProvider(t *testing.T) {

	tests := []struct {
		name string
		sp   *SearchParam
		want string
	}{
		{
			name: "empty-provider",
			sp:   &SearchParam{},
			want: "",
		},
		{
			name: "provider-private",
			sp:   &SearchParam{Provider: Private},
			want: "/anbieter:" + string(Private),
		},
		{
			name: "provider-commercial",
			sp:   &SearchParam{Provider: Commercial},
			want: "/anbieter:" + string(Commercial),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.sp.fmtProvider(); got != tt.want {
				t.Errorf("SearchParam.fmtProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fmtOfferType(t *testing.T) {

	tests := []struct {
		name string
		sp   *SearchParam
		want string
	}{
		{
			name: "empty-offer-type",
			sp:   &SearchParam{},
			want: "",
		},
		{
			name: "offer-type-offer",
			sp:   &SearchParam{OfferType: Offer},
			want: "/anzeige:" + string(Offer),
		},
		{
			name: "offer-type-wanted",
			sp:   &SearchParam{OfferType: Wanted},
			want: "/anzeige:" + string(Wanted),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.sp.fmtOfferType(); got != tt.want {
				t.Errorf("SearchParam.fmtOfferType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fmtLocationID(t *testing.T) {

	tests := []struct {
		name string
		sp   *SearchParam
		want string
	}{
		{
			name: "empty-location",
			sp:   &SearchParam{},
			want: "",
		},
		{
			name: "location-berlin",
			sp:   &SearchParam{Location: "3331"},
			want: "l3331",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.sp.fmtLocationID(); got != tt.want {
				t.Errorf("SearchParam.fmtLocationID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fmtRadius(t *testing.T) {

	tests := []struct {
		name string
		sp   *SearchParam
		want string
	}{
		{
			name: "empty-radius",
			sp:   &SearchParam{},
			want: "",
		},
		{
			name: "radius-whole-location",
			sp:   &SearchParam{Radius: WholePlace},
			want: "",
		},
		{
			name: "radius-fifty",
			sp:   &SearchParam{Radius: FiftyKM},
			want: "r" + string(FiftyKM),
		},
		{
			name: "radius-one-hundred-fifty",
			sp:   &SearchParam{Radius: OneHundredFiftyKM},
			want: "r" + string(OneHundredFiftyKM),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.sp.fmtRadius(); got != tt.want {
				t.Errorf("SearchParam.fmtRadius() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fmtSpecificParameter(t *testing.T) {

	tests := []struct {
		name string
		sp   *SearchParam
		want string
	}{
		{
			name: "empty-parameters",
			sp:   &SearchParam{},
			want: "",
		},
		{
			name: "parameter-car",
			sp: &SearchParam{
				SpecificParameter: SpecificParameter{
					CarManufacturer: "bmw",
				},
			},
			want: "+" + string(CarManufacturer) + ":bmw",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sp.fmtSpecificParameter(); got != tt.want {
				t.Errorf("SearchParam.fmtSpecificParameter() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_toURL(t *testing.T) {

	tests := []struct {
		name string
		sp   *SearchParam
		want string
	}{
		{
			name: "empty-params",
			sp:   &SearchParam{},
			want: baseURL + "/seite:1",
		},
		{
			name: "param-car",
			sp: &SearchParam{
				Category: Cars,
			},
			want: baseURL + "/seite:1/c" + string(Cars),
		},
		{
			name: "param-page",
			sp: &SearchParam{
				Page: 23,
			},
			want: baseURL + "/seite:23",
		},
		{
			name: "param-price",
			sp: &SearchParam{
				PriceFrom: 1000,
				PriceTo:   5000,
			},
			want: baseURL + "/preis:1000:5000/seite:1",
		},
		{
			name: "param-basic",
			sp: &SearchParam{
				Category:  Cars,
				OfferType: Offer,
				Location:  "3331",
				Radius:    TenKM,
				SpecificParameter: SpecificParameter{
					CarManufacturer: "bmw",
				},
			},
			want: baseURL + "/anzeige:angebote/seite:1/c216l3331r10+autos.marke_s:bmw",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sp.toURL(); got != tt.want {
				t.Errorf("SearchParam.toURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
