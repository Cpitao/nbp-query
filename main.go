package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
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

type NbpARate struct {
	Table    string  `json:"table"`
	Currency string  `json:"currency"`
	Code     string  `json:"code"`
	Rate     []ARate `json:"rates"`
}

func getNbpARate() (NbpARate, error) {
	queryUrl := "http://api.nbp.pl/api/exchangerates/rates/a/gbp/?format=json"
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
func (v Currency) convert() (Currency, float64) {
	if v.Name == "gbp" || v.Name == "pln" {
		nbpRate, err := getNbpARate()
		if err != nil {
			return Currency{}, -1
		}

		// return values rounded to 2 decimal places
		if v.Name == "pln" {
			return Currency{"gbp", math.Floor(100*v.Value/nbpRate.Rate[0].Mid) / 100}, nbpRate.Rate[0].Mid
		} else if v.Name == "gbp" {
			return Currency{"pln", math.Floor(100*v.Value*nbpRate.Rate[0].Mid) / 100}, nbpRate.Rate[0].Mid
		}
	}
	return Currency{}, -1
}

func main() {
	tmpl := template.Must(template.ParseFiles("template.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Gbps    string
			Plns    string
			IsRate  bool
			Rate    string
			IsError bool
			Error   string
		}{"", "", false, "", false, ""}

		gbpString := r.FormValue("gbp")
		plnString := r.FormValue("pln")

		if gbpString == "" && plnString == "" {
			tmpl.Execute(w, data)
			return
		}

		if gbpString != "" && plnString != "" {
			data.IsError = true
			data.Error = "One field must remain empty"
			tmpl.Execute(w, data)
			return
		}

		var currencyType string
		var value float64
		var err error

		if gbpString != "" {
			currencyType = "gbp"
			value, err = strconv.ParseFloat(gbpString, 64)
		} else {
			currencyType = "pln"
			value, err = strconv.ParseFloat(plnString, 64)
		}

		if err != nil || value < 0 {
			data.IsError = true
			data.Error = "Invalid value"
			tmpl.Execute(w, data)
			return
		}

		value = math.Floor(100*value) / 100
		inputCurrency := Currency{currencyType, value}
		outputCurrency, rate := inputCurrency.convert()

		if rate < 0 {
			data.IsError = true
			data.Error = "Unable to convert values"
			tmpl.Execute(w, data)
			return
		}

		if inputCurrency.Name == "gbp" {
			data.Gbps = fmt.Sprintf("%.2f", inputCurrency.Value)
			data.Plns = fmt.Sprintf("%.2f", outputCurrency.Value)
		} else {
			data.Gbps = fmt.Sprintf("%.2f", outputCurrency.Value)
			data.Plns = fmt.Sprintf("%.2f", inputCurrency.Value)
		}

		data.IsRate = true
		data.Rate = fmt.Sprintf("%.2f", rate)

		tmpl.Execute(w, data)
	})

	http.ListenAndServe(":80", nil)
}
