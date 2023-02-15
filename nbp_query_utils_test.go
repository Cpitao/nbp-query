package application

import (
	"math"
	"testing"
)

func currencyEqualWithTolerance(a, b Currency, epsilon float64) bool {
	return math.Abs(a.Value-b.Value) < epsilon && a.Name == b.Name
}

type convertTest struct {
	currency Currency
	target   string
	rate     float64
	expected Currency
}

var convertTests = []convertTest{
	{Currency{"pln", 10.0}, "gbp", 5.0, Currency{"gbp", 2.0}},
	{Currency{"gbp", 10.0}, "gbp", 5.0, Currency{"pln", 50.0}},
	{Currency{"usd", 10.0}, "usd", 5.1234, Currency{"pln", 51.23}},
	{Currency{"pln", 10.0}, "usd", 5.1234, Currency{"usd", 1.95}}}

func TestConvert(t *testing.T) {
	var epsilon = 1e-10
	for _, test := range convertTests {
		out := test.currency.convert(test.target, test.rate)
		if !currencyEqualWithTolerance(out, test.expected, epsilon) {
			t.Errorf("Output Currency{\"%s\", %f} not equal Currency{\"%s\", %f}",
				out.Name, out.Value, test.expected.Name, test.expected.Value)
		}
	}
}

type VerifyCurrencyCodeTest struct {
	code     string
	expected bool
}

var verifyCurrencyCodeTests = []VerifyCurrencyCodeTest{
	{"ABC", true},
	{"aBc", true},
	{"abc", true},
	{"abcd", false},
	{"abc1", false},
	{"ab1", false},
	{"123", false}}

func TestVerifyCurrencyCode(t *testing.T) {
	for _, test := range verifyCurrencyCodeTests {
		out := verifyCurrencyCode(test.code)
		if out != test.expected {
			t.Errorf("Output %t not equal %t for %s", out, test.expected, test.code)
		}
	}
}
