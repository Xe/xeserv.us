package xonotic

import "testing"

func TestGetStatus(t *testing.T) {
	c := Dial("10.0.0.18", "26000")

	stat, err := c.Status()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%#v", stat)
}
