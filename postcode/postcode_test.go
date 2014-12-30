package postcode

import (
	"testing"
)

func TestXSNews(t *testing.T) {
	entity, e := Get("1175RD", "9", "")
	if e != nil {
		t.Error("Failed calling webservice: " + e.Error());
	}

	if entity.Valid == false {
		t.Error("Zipcode resolving failed");
	}
	if entity.Municipality != "Haarlemmermeer" {
		t.Error("Wrong municipality: " + entity.Municipality)
	}
}
