package vies

import (
	"testing"
)

func TestXSNews(t *testing.T) {
	entity, e := Get("NL815836946B01")
	if e != nil {
		t.Error("Failed calling webservice: " + e.Error());
	}

	if entity.Valid != true {
		t.Error("Invalid VAT number")
	}
	if entity.Name != "XS NEWS B.V." {
		t.Error("Invalid company received: " + entity.Name)
	}
}
