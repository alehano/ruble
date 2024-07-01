// Official rates of russian ruble for a specific date
// http://www.cbr.ru/scripts/Root.asp?PrtId=SXML
package ruble

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	"strings"

	"golang.org/x/net/html/charset"
)

type valCurs struct {
	XMLName struct{}   `xml:"ValCurs"`
	Valutes []Currency `xml:"Valute"`
}

type Currency struct {
	Code    string  `xml:"CharCode"`
	Nominal int     `xml:"Nominal"`
	Name    string  `xml:"Name"`
	Value   float64 `xml:"Value"`
}

func GetCurrenciesByCode(date time.Time) (map[string]float64, error) {
	curs, err := GetCurrencies(date)
	if err != nil {
		return nil, err
	}
	return CurrenciesByCode(curs), nil
}

func GetCurrencies(date time.Time) ([]Currency, error) {
	valCurs := valCurs{}
	res := []Currency{}
	data, err := getXML(date)
	if err != nil {
		return res, err
	}

	// replace , to .
	data = []byte(strings.Replace(string(data), ",", ".", -1))

	buf := bytes.NewBuffer(data)
	decoder := xml.NewDecoder(buf)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&valCurs)
	return valCurs.Valutes, err
}

func CurrenciesByCode(curs []Currency) map[string]float64 {
	res := map[string]float64{}
	for _, cur := range curs {
		res[cur.Code] = RoundPlus(cur.Value*float64(cur.Nominal), 2)
	}
	return res
}

func getXML(date time.Time) ([]byte, error) {
	const dateForm = "02/01/2006"
	url := fmt.Sprintf("https://www.cbr.ru/scripts/XML_daily.asp?date_req=%s", date.Format(dateForm))
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.75 Safari/537.36")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	resp, err := client.Do(req)

	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err

}

func Round(x float64) float64 {
	v, frac := math.Modf(x)
	if x > 0.0 {
		if frac > 0.5 || (frac == 0.5 && uint64(v)%2 != 0) {
			v += 1.0
		}
	} else {
		if frac < -0.5 || (frac == -0.5 && uint64(v)%2 != 0) {
			v -= 1.0
		}
	}
	return v
}

func RoundPlus(x float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return Round(x*shift) / shift
}
