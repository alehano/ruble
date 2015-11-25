# ruble
Official rates of russian ruble for a specific date.

```go
curs, err := GetCurrenciesByCode(time.Now())

// How many rubles in 100 USD?
fmt.Printf("A hundred bucks = %.2f rubles\n", 100*curs["USD"])
```
