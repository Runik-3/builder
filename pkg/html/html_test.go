package html

import (
	"testing"
)

func TestParseFirstPara(t *testing.T) {
	want := "this is a test"

	content, _ := (ParseFirstPara("<p><b>this</b> is a test</p>"))
	if content != want {
		t.Errorf("got %s want %s", content, want)
	}

	content, _ = (ParseFirstPara("<p>this is a <b>test</b></p><p>asdasd</p>"))
	if content != want {
		t.Errorf("got %s want %s", content, want)
	}

	content, _ = (ParseFirstPara("<p>asdasd</p><p>this is a <b>test</b></p>"))
	if content != want {
		t.Errorf("got %s want %s", content, want)
	}

	_, err := (ParseFirstPara("<p>this is a test</p><p>asdasd</p>"))
	if err == nil {
		t.Errorf("expected error")
	}
}
