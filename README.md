# nbp-query
Application allows converting from and to PLN from all common currencies, using exchange rates from [NBP API](http://api.nbp.pl).<br>
Currency type format is 3 letters as specified in [ISO 4217](https://pl.wikipedia.org/wiki/ISO_4217) (some of the supported formats include EUR, CHF, USD, etc.).

### Usage
Clone this repository.

From the command line, within the directory containing *main.go* file run `go run main.go`.

The page will be served on localhost at port 80.

### Tests
To run unit tests, use `go test nbp_query/nbp`.
