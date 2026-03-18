package wikitext

import _ "embed"

//go:embed fixtures/test_wikitext_lg_1.txt
var sample_wikitext_lg_1 string
//go:embed fixtures/test_wikitext_lg_2.txt
var sample_wikitext_lg_2 string
//go:embed fixtures/test_wikitext_lg_3.txt
var sample_wikitext_lg_3 string
//go:embed fixtures/test_wikitext_lg_4.txt
var sample_wikitext_lg_4 string

//go:embed fixtures/test_wikitext_sm.txt
var sample_wikitext_sm string

var benchmark_sample_text = []string{
	sample_wikitext_sm,
	sample_wikitext_lg_1,
	sample_wikitext_lg_2,
	sample_wikitext_lg_3,
	sample_wikitext_lg_4,
}
