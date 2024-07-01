# ruble

Official CBRF currency rates of Russian ruble for a specific date.
Source: https://www.cbr.ru

```go
curs, err := GetCurrenciesByCode(time.Now())

// How many rubles in 100 USD?
fmt.Printf("A hundred bucks = %.2f rubles\n", 100*curs["USD"])
```
