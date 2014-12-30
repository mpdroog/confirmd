package openkvk

import (
	"testing"
)

func TestXSNews(t *testing.T) {
	entity, e := Get("08145640")
	if e != nil {
		t.Error("Failed calling webservice: " + e.Error());
	}
	if entity.Legal != "B.V." {
		t.Error("Entity wrong legal form: " + entity.Legal)
	}
	if entity.Name != "XS News B.V." {
		t.Error("Entity wrong name: " + entity.Name);
	}
}

func TestIts(t *testing.T) {
	entity, e := Get("37113555")
	if e != nil {
		t.Error("Failed calling webservice: " + e.Error());
	}
	if entity.Name != "Its Hosted" {
		t.Error("Invalid name found for entity: " + entity.Name)
	}
}
