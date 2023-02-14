package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	tmpl := template.Must(template.ParseFiles("template.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			CurrencyType  string
			CurrencyCount string
			Plns          string
			IsRate        bool
			Rate          string
			IsError       bool
			Error         string
		}{"GBP", "", "", false, "", false, ""} // default values

		currencyTypeString := strings.ToLower(r.FormValue("type"))
		currencyRegex, _ := regexp.Compile("^[a-z]{3}$")
		if !currencyRegex.MatchString(currencyTypeString) {
			data.IsError = true
			data.Error = "Invalid currency type format"
			tmpl.Execute(w, data)
			return
		}

		inValueString := r.FormValue("other")
		plnString := r.FormValue("pln")
		fmt.Printf("%s %s %s", currencyTypeString, inValueString, plnString)

		if inValueString == "" && plnString == "" {
			tmpl.Execute(w, data)
			return
		}

		if inValueString != "" && plnString != "" { // handled client-side in the first place
			data.IsError = true
			data.Error = "One field must remain empty"
			tmpl.Execute(w, data)
			return
		}

		var currencyType string
		var value float64
		var err error

		if inValueString != "" {
			currencyType = currencyTypeString
			value, err = strconv.ParseFloat(inValueString, 64)
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

		value = math.Round(100*value) / 100
		inputCurrency := Currency{currencyType, value}
		outputCurrency, rate := inputCurrency.convert(currencyTypeString)

		if rate < 0 {
			data.IsError = true
			data.Error = "Unable to convert values"
			tmpl.Execute(w, data)
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
		data.Rate = fmt.Sprintf("%.2f", rate)
		data.CurrencyType = strings.ToUpper(currencyTypeString)

		tmpl.Execute(w, data)
	})

	http.ListenAndServe(":80", nil)
}
