package openkvk
/**
 * openkvk.nl - Chamber of Commerce API for NL
 */
import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

const URL = `http://officieel.openkvk.nl/json/`

type rawEntity struct {
	Name         string `json:"rechtspersoon"`
	Area         string `json:"vestigingsnummer"`
	Address      string `json:"adres"`
	City         string `json:"woonplaats"`
	CoC          string `json:"kvk"`
	TradeNames   map[string][]string `json:"handelsnamen"`
	Zipcode      string `json:"postcode"`
	Type         string `json:"type"`
	CoCs         string `json:"kvks"`
}
type KvkEntity struct {
	Name     string
	CoC      string
	Legal    string
	Area     string

	Address  string
	City     string
	TradeNames []string
	Zipcode  string
}

var (
	bv string = "Rechtspersoon"
)

func getKvK(number string) ([]rawEntity, error) {
	res, e := http.Get(URL + number)
	if e != nil {
		return nil, e
	}
	b, e := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if e != nil {
		return nil, e
	}

	var c []rawEntity
	if e := json.Unmarshal(b, &c); e != nil {
		return nil, e
	}
	return c, nil
}

func Get(number string) (*KvkEntity, error) {
	k, e := getKvK(number)
	if e != nil {
		return nil, e
	}

	var (
		entity KvkEntity
	)
	for _, item := range k {
		if (item.Type == bv) {
			entity.Legal = "B.V."
		} else {
			entity.CoC = item.CoCs
			entity.Name = item.Name
			entity.Area = item.Area
			entity.Address = item.Address
			entity.City = item.City
			entity.TradeNames = item.TradeNames["bestaand"]
			entity.Zipcode = item.Zipcode
		}
	}
	return &entity, nil
}
