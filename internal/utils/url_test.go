package utils

import (
	"testing"
)

func TestUrlValidation(t *testing.T) {
	cases := [][]string{
		{"https://malazan.fandom.com", "https://malazan.fandom.com/api.php"},
		{"https://malazan.fandom.com/api.php", "https://malazan.fandom.com/api.php"},
		// handle url with lang code
		{"https://amanecer-rojo.fandom.com/es/wiki/Amanecer_Rojo_Wiki", "https://amanecer-rojo.fandom.com/es/api.php"},
		// invalid language code
		{"https://malazan.fandom.com/xx/wiki/front_page", "https://malazan.fandom.com/api.php"},
	}

	for _, c := range cases {
		r, _ := NormalizeUrl(c[0])
		if r != c[1] {
			t.Fatalf("\nWant:\n%s\n\nRecieved:\n%s", c[1], r)
		}
	}
}
