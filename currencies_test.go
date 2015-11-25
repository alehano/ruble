package currencies

import (
	"testing"
	"time"

	"fmt"
)

func TestCur(t *testing.T) {
	data, err := GetCurrencies(time.Now())
	if err != nil {
		t.Error(err)
	}
	fmt.Println(data)
	if len(data) == 0 {
		t.Error("Empty result")
	}
}

func TestCurByCode(t *testing.T) {
	data, err := GetCurrenciesByCode(time.Now())
	if err != nil {
		t.Error(err)
	}
	fmt.Println(data)
	if len(data) == 0 {
		t.Error("Empty result")
	}
}
