package wikitext

import (
	"reflect"
	"testing"
)

type TknCase struct {
	input    string
	expected []Token
}

func TestTokenizer(t *testing.T) {
	cases := []TknCase{
		//	{input: "'''Akari au Raa''' was a [[Gold]], the progenitor of [[House Raa]], and one of the founders of the [[The Society|Society]].", expected: []Token{{Type: "text", Value: "Akari au Raa was a "}, {Type: "link", Value: "Gold, "}, {Type: "text", Value: "the progenitor of "}, {Type: "link", Value: "House Raa, "}, {Type: "text", Value: "and one of the founders of the "}, {Type: "link", Value: "The Society|Society. "}}},
		//	{input: "'''Achlys-9 '''is a deadly poison gas. It is used in executions and the quarantine of mines. [[Bryn]] was killed with this.\n[[hu:Akhlüsz-9]]\n[[es:Aclis-9]]\n[[Category:Materials]]", expected: []Token{{Type: "text", Value: "Achlys-9 is a deadly poison gas. It is used in executions and the quarantine of mines. "}, {Type: "link", Value: "Bryn "}, {Type: "text", Value: "was killed with this. "}, {Type: "link", Value: "hu:Akhlüsz-9 "}, {Type: "link", Value: "es:Aclis-9 "}, {Type: "link", Value: "Category:Materials "}}},
		//	{input: "#REDIRECT [[The Jackal]]", expected: []Token{{Type: "text", Value: "#REDIRECT "}, {Type: "link", Value: "The Jackal "}}}, // handle redirects gracefully
		//	{input: "[[File:Agea During Terraforming.jpg|thumb|Agea during the terraforming of Mars]]\n<p class=\"MsoNormal\">'''Agea''' is the capital city of [[Mars]]. It is located in the [[Valles Marineris]], the largest canyon in the Solar System.  </p>\n\n<p class=\"MsoNormal\">The citadel of the [[ArchGovernor]], [[Nero au Augustus]], is located in Agea, as is the [[House Bellona]] family estate. Nero holds court within the city's Forum.</p>\n\n<p class=\"MsoNormal\">Described as beautiful but strange, life in the city is both fast and cheap. It is well known for its nightlife, with rooftop nightclubs, [[gravMixers]] and [[NoiseBubbles]].</p>\n\n<p class=\"MsoNormal\">Agea is home to the [[Agea Martial Club]], where [[Gold]]'s are able to duel and practice their fighting skills against their peers.</p>\n\n<p class=\"MsoNormal\">Once [[The Institute]] trials are complete, a grand festival is held in Agea to honour the newly graduated [[Peerless Scarred]].</p>\n\n== Gallery ==\n<gallery>\nFile:Agea mars 720 pce.jpg|Agea, Mars 720 PCE\nFile:Agea Dueling Club.jpg|A dueling club in Agea, 725 PCE\nFile:725 Agea Skyline at Night.jpg|The Agea skyline at night, 725 PCE\nFile:Agea Forum.jpg|Agea's Forum, 725 PCE\n</gallery>\n[[es:Agea]]\n[[hu:Égea]]\n[[Category:Locations]]\n[[Category:Cities]]", expected: []Token{{Type: "link", Value: "File:Agea During Terraforming.jpg|thumb|Agea during the terraforming of Mars "}, {Type: "text", Value: "Agea is the capital city of "}, {Type: "link", Value: "Mars. "}, {Type: "text", Value: "It is located in the "}, {Type: "link", Value: "Valles Marineris, "}, {Type: "text", Value: "the largest canyon in the Solar System. The citadel of the "}, {Type: "link", Value: "ArchGovernor, "}, {Type: "link", Value: "Nero au Augustus, "}, {Type: "text", Value: "is located in Agea, as is the "}, {Type: "link", Value: "House Bellona "}, {Type: "text", Value: "family estate. Nero holds court within the city's Forum. Described as beautiful but strange, life in the city is both fast and cheap. It is well known for its nightlife, with rooftop nightclubs, "}, {Type: "link", Value: "gravMixers "}, {Type: "text", Value: "and "}, {Type: "link", Value: "NoiseBubbles. "}, {Type: "text", Value: "Agea is home to the "}, {Type: "link", Value: "Agea Martial Club, "}, {Type: "text", Value: "where "}, {Type: "link", Value: "Gold's "}, {Type: "text", Value: "are able to duel and practice their fighting skills against their peers. Once "}, {Type: "link", Value: "The Institute "}, {Type: "text", Value: "trials are complete, a grand festival is held in Agea to honour the newly graduated "}, {Type: "link", Value: "Peerless Scarred. "}, {Type: "text", Value: "== Gallery == File:Agea mars 720 pce.jpg|Agea, Mars 720 PCE File:Agea Dueling Club.jpg|A dueling club in Agea, 725 PCE File:725 Agea Skyline at Night.jpg|The Agea skyline at night, 725 PCE File:Agea Forum.jpg|Agea's Forum, 725 PCE "}, {Type: "link", Value: "es:Agea "}, {Type: "link", Value: "hu:Égea "}, {Type: "link", Value: "Category:Locations "}, {Type: "link", Value: "Category:Cities "}}},
		//	{input: "{{Gold_Character\n|title1 =  [[File:Gold_Sigil.png|left|50px]] Ajax [[File:Gold_Sigil.png|right|50px]]\n|image1 =  Ajax.png \n|caption1 = \n|alias(es) = Storm Knight (as of 754 PCE)\n|gender = Male|age = ~23 (Dark Age)\n}}\n'''Ajax au Grimmus''' is the son of [[Aja au Grimmus]] and [[Atlas au Raa]].", expected: []Token{{Type: "table", Value: "{{Gold_Character |title1 = [[File:Gold_Sigil.png|left|50px]] Ajax [[File:Gold_Sigil.png|right|50px]] |image1 = Ajax.png |caption1 = |alias(es) = Storm Knight (as of 754 PCE) |gender = Male|age = ~23 (Dark Age) }} "}, {Type: "text", Value: "Ajax au Grimmus is the son of "}, {Type: "link", Value: "Aja au Grimmus "}, {Type: "text", Value: "and "}, {Type: "link", Value: "Atlas au Raa. "}}},
		{input: "'''Denna''' {{pron|PR|/'dɛne/}}<ref>https://youtu.be/MPEB6NAGoYk?t=1041</ref> is the primary female figure in ''[[The Name of the Wind]]''; she is arguably the main romantic interest of Kvothe, who holds an irresistible fascination with her.", expected: []Token{}},
	}

	for _, c := range cases {
		r := tokenizer(c.input)

		if !reflect.DeepEqual(r, c.expected) {
			t.Fatalf("Want:\n%v\n\nRecieved:\n%s", c.expected, r)
		}
	}
}

func TestCleanText(t *testing.T) {
	cases := [][]string{
		{"<p>This is a test string</p>", "This is a test string"},
		{"<img src=\"source.com\">", ""},
		{"'''Welcome''' <img src=\"source.com\">hi", "Welcome hi"},
		{"<div class=\"top-styling\">This is a <span>test</span> string</div>", "This is a test string"},
		{"'''''Welcome''''' hi", "Welcome hi"},
	}

	for _, c := range cases {
		r := cleanText(c[0])
		if r != c[1] {
			t.Fatalf("Want:\n%s\n\nRecieved:\n%s", c[1], r)
		}
	}
}
