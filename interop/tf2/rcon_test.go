package tf2

import (
	"log"
	"testing"
)

func TestGetStatus(t *testing.T) {
	s, err := Query("10.0.0.5:27025", "cqcontrol")
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%#v", s)
}
