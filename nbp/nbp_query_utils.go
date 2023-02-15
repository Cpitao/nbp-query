package nbp

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strings"
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

func QueryNbpRate(target string) (NbpARate, error) {
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

func (c Currency) convert(target string, rate float64) Currency {
	if c.Name == "pln" {
		return Currency{target, math.Round(100*c.Value/rate) / 100}
	} else {
		return Currency{"pln", math.Round(100*c.Value*rate) / 100}
	}
}

func ExchangeCurrency(c Currency, target string) (Currency, float64) {
	rate, err := QueryNbpRate(target)
	if err != nil {
		return Currency{}, -1
	}

	return c.convert(target, rate.Rate[0].Mid), rate.Rate[0].Mid
}

func VerifyCurrencyCode(code string) bool {
	code = strings.ToLower(code)
	currencyRegex, _ := regexp.Compile("^[a-z]{3}$")
	return currencyRegex.MatchString(code)
}
