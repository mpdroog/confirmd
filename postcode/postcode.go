package postcode
// Postal code lookup for The Netherlands
// using public DB available at postcode.nl

import (
	"fmt"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"errors"
	"confirmd/config"
)

const URL = `https://api.postcode.nl/rest/addresses/%s/%s/%s`

type Entity struct {
	Valid               bool
	Street              string
	HouseNumber         int
	HouseNumberAddition *string
	PostCode            string
	City                string
	Municipality        string
	Province            string
	AddressType         string
	Purposes            []string
	HouseNumerAdditions []string
}

func Get(zip string, houseNo string, houseNoAdd string) (Entity, error) {
	client := &http.Client{}
	req, e := http.NewRequest("POST", fmt.Sprintf(
			URL, url.QueryEscape(zip), url.QueryEscape(houseNo),
			url.QueryEscape(houseNoAdd),
	), nil)
	if e != nil {
		return Entity{}, e
	}
	req.SetBasicAuth(config.Pref.Postcode.Username, config.Pref.Postcode.Password)

	res, e := client.Do(req)
	if e != nil {
		return Entity{}, e
	}

	if res.StatusCode == 404 {
		return Entity{}, nil
	}
	if res.StatusCode != 200 {
		return Entity{}, errors.New("Weird HTTP-error " + res.Status)
	}

	b, e := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if e != nil {
		return Entity{}, e
	}

	var entity Entity
	if e := json.Unmarshal(b, &entity); e != nil {
		return Entity{}, e
	}
	entity.Valid = true
	return entity, nil
}
