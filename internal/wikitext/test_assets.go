package wikitext

import _ "embed"

//go:embed fixtures/test_wikitext_lg.txt
var sample_wikitext_lg string

//go:embed fixtures/test_wikitext_sm.txt
var sample_wikitext_sm string
