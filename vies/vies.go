package vies
/**
 * VIES EU VAT number validation
 * @link http://ec.europa.eu/taxation_customs/vies/?locale=en
 */

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"io"
	"strings"
	"encoding/xml"
	"errors"
)

const soapTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/" encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <Header/>
    <Body>
        <checkVat xmlns="urn:ec.europa.eu:taxud:vies:services:checkVat:types">
            <countryCode>%s</countryCode>
            <vatNumber>%s</vatNumber>
        </checkVat>
    </Body>
</Envelope>`

const URL = `http://ec.europa.eu/taxation_customs/vies/services/checkVatService`

type Entity struct {
	CountryCode string `xml:"countryCode"`
	VatNumber   string `xml:"vatNumber"`
	RequestDate string `xml:"requestDate"`
	Valid       bool   `xml:"valid"`
	Name        string `xml:"name"`
	Address     string `xml:"address"`
}

// Use hardcoded string for request
func soapify(iso2 string, vat string) string {
	return fmt.Sprintf(soapTemplate, iso2, vat)
}

// Parse raw string for error
// if no error convert to data structure
func parse(input string) (*Entity, error) {
	b := strings.Index(input, "<faultstring>")
	if b != -1 {
		// Error
		e := strings.Index(input, "</faultstring>")
		if e == -1 {
			return nil, errors.New("XML Parse error: Missing </faultstring>")
		}
		return nil, errors.New("SOAP Error: " + input[b:e] )
	}

	entity := new(Entity)
	decoder := xml.NewDecoder(strings.NewReader(input))
	for {
		t, e := decoder.Token()
		if e != io.EOF && e != nil {
			return nil, e
		}
		if t == nil {
			break;
		}

		elem, found := t.(xml.StartElement)
		if found && elem.Name.Local == "checkVatResponse" {
			if e := decoder.DecodeElement(&entity, &elem); e != nil {
				return nil, e
			}
		}
	}

	return entity, nil
}

func Get(number string) (Entity, error) {
	body := soapify(number[:2], number[2:])

	res, e := http.Post(URL, "text/xml", strings.NewReader(body))
	if e != nil {
		return Entity{}, e
	}
	b, e := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if e != nil {
		return Entity{}, e
	}

	entity, e := parse(string(b))
	if e != nil {
		return Entity{}, e
	}
	return *entity, nil
}
