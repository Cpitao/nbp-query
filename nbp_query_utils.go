package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

type Currency struct {
	Name  string
	Value float64
}

type ARate struct {
	No            string  `json:"no"`
	EffectiveDate string  `json:"effectiveDate"`
	Mid           float64 `json:"mid"`
}

type NbpARate struct { // structure to hold table A type of data
	Table    string  `json:"table"`
	Currency string  `json:"currency"`
	Code     string  `json:"code"`
	Rate     []ARate `json:"rates"`
}

func getNbpARate(target string) (NbpARate, error) {
	queryUrl := fmt.Sprintf("http://api.nbp.pl/api/exchangerates/rates/a/%s/?format=json", target)
	resp, err := http.Get(queryUrl)
	if err != nil {
		return NbpARate{}, err
	}
	defer resp.Body.Close()
	if err != nil {
		return NbpARate{}, err
	}

	var nbpRate NbpARate
	err = json.NewDecoder(resp.Body).Decode(&nbpRate)
	if err != nil {
		return NbpARate{}, err
	}
	return nbpRate, nil
}

/* return value of returned currency and conversion rate */
func (v Currency) convert(target string) (Currency, float64) {
	nbpRate, err := getNbpARate(target)
	if err != nil {
		return Currency{}, -1
	}

	// return values rounded to 2 decimal places
	if v.Name == "pln" {
		return Currency{target, math.Round(100*v.Value/nbpRate.Rate[0].Mid) / 100}, nbpRate.Rate[0].Mid
	} else {
		return Currency{"pln", math.Round(100*v.Value*nbpRate.Rate[0].Mid) / 100}, nbpRate.Rate[0].Mid
	}
}
