# ruble
Official currencies of russian ruble for specific date

```go
curs, err := GetCurrenciesByCode(time.Now())

// How many rubles in 100 USD?
fmt.Printf("A hundred bucks = %.2f rubles\n", curs["USD"]*100)
```
