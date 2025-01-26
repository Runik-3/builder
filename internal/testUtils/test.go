package testUtils

import (
	"strings"
	"testing"
)

func IsEqual(t *testing.T, val1 any, val2 any, message string) {
	if val1 != val2 {
		if message != "" {
			t.Fatal(message)
		} else {
			t.Fatalf("%v not equal to %v", val1, val2)
		}
	}
}

func Contains(t *testing.T, str string, subStr string) {
	if !strings.Contains(str, subStr) {
		t.Fatalf("%v does not contain %v", str, subStr)
	}
}
