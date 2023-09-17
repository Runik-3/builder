package html

import (
	"testing"
)

func TestParseFirstPara(t *testing.T) {
	string, _ := (ParseFirstPara("<p><b>this</b> is a test</p>"))
	want := "this is a test"
	if string != want {
		t.Errorf("got %s want %s", string, want)
	}

	_, err := (ParseFirstPara("<p>asd <b>this</b> is a test</p>"))
	if err == nil {
		t.Errorf("expected error")
	}
}
