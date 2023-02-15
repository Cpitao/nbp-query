package main

import (
	"fmt"
	"html/template"
	"math"
	nbp "nbp_query/nbp"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	tmpl := template.Must(template.ParseFiles("template.html"))

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		data := struct {
			CurrencyType  string
			CurrencyCount string
			Plns          string
			IsRate        bool
			Rate          string
			IsError       bool
			Error         string
		}{"GBP", "", "", false, "", false, ""} // default values

		otherCurrencyValue := request.FormValue("other")
		plnValue := request.FormValue("pln")

		if otherCurrencyValue == "" && plnValue == "" {
			tmpl.Execute(writer, data)
			return
		}

		currencyCode := request.FormValue("type")
		if !nbp.VerifyCurrencyCode(currencyCode) {
			data.IsError = true
			data.Error = "Invalid currency type format"
			tmpl.Execute(writer, data)
			return
		}

		if otherCurrencyValue != "" && plnValue != "" {
			data.IsError = true
			data.Error = "One field must remain empty"
			tmpl.Execute(writer, data)
			return
		}

		var currencyType string
		var value float64
		var err error

		if otherCurrencyValue != "" {
			currencyType = currencyCode
			value, err = strconv.ParseFloat(otherCurrencyValue, 64)
		} else {
			currencyType = "pln"
			value, err = strconv.ParseFloat(plnValue, 64)
		}

		if err != nil || value < 0 {
			data.IsError = true
			data.Error = "Invalid value"
			tmpl.Execute(writer, data)
			return
		}

		value = math.Round(100*value) / 100
		inputCurrency := nbp.Currency{Name: currencyType, Value: value}
		outputCurrency, rate := nbp.ExchangeCurrency(inputCurrency, currencyCode)

		if rate < 0 {
			data.IsError = true
			data.Error = "Unable to convert values"
			tmpl.Execute(writer, data)
			return
		}

		if inputCurrency.Name != "pln" {
			data.CurrencyCount = fmt.Sprintf("%.2f", inputCurrency.Value)
			data.Plns = fmt.Sprintf("%.2f", outputCurrency.Value)
		} else {
			data.CurrencyCount = fmt.Sprintf("%.2f", outputCurrency.Value)
			data.Plns = fmt.Sprintf("%.2f", inputCurrency.Value)
		}

		data.IsRate = true
		data.Rate = fmt.Sprintf("%.4f", rate)
		data.CurrencyType = strings.ToUpper(currencyCode)

		tmpl.Execute(writer, data)
	})

	http.ListenAndServe(":80", nil)
}
