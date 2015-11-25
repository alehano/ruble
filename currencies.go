// Official currencies of russian ruble for specific date
// http://www.cbr.ru/scripts/Root.asp?PrtId=SXML
package currencies

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
	res := map[string]float64{}
	curs, err := GetCurrencies(date)
	if err != nil {
		return res, err
	}
	for _, cur := range curs {
		res[cur.Code] = RoundPlus(cur.Value*float64(cur.Nominal), 2)
	}
	return res, nil
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

func getXML(date time.Time) ([]byte, error) {
	const dateForm = "02/01/2006"
	url := fmt.Sprintf("http://www.cbr.ru/scripts/XML_daily.asp?date_req=%s", date.Format(dateForm))
	resp, err := http.Get(url)
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
