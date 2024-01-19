package utils

import "testing"

func TestUrlValidation(t *testing.T) {
	cases := [][]string{
		{"https://malazan.fandom.com", "https://malazan.fandom.com/api.php"},
		{"https://malazan.fandom.com/api.php", "https://malazan.fandom.com/api.php"},
	}

	for _, c := range cases {
		r, _ := FormatUrl(c[0])
		if r != c[1] {
			t.Fatalf("\nWant:\n%s\n\nRecieved:\n%s", c[1], r)
		}
	}
}
