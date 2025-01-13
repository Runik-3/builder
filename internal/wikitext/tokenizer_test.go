package wikitext

import (
	_ "embed"
	"testing"
)

func TestTokenizer(t *testing.T) {
	cases := [][]string{
		{"'''Akari au Raa''' was a [[Gold]], the progenitor of [[House Raa]], and one of the founders of the [[The Society|Society]].", `[{"Type":"text","Value":"Akari au Raa was a "},{"Type":"link","Value":"Gold"},{"Type":"text","Value":", the progenitor of "},{"Type":"link","Value":"House Raa"},{"Type":"text","Value":", and one of the founders of the "},{"Type":"link","Value":"The Society|Society"},{"Type":"text","Value":"."}]`},
		{"'''Achlys-9 '''is a deadly poison gas. It is used in executions and the quarantine of mines. [[Bryn]] was killed with this.\n[[hu:Akhlüsz-9]]\n[[es:Aclis-9]]\n[[Category:Materials]]", `[{"Type":"text","Value":"Achlys-9 is a deadly poison gas. It is used in executions and the quarantine of mines. "},{"Type":"link","Value":"Bryn"},{"Type":"text","Value":" was killed with this.\n"},{"Type":"link","Value":"hu:Akhlüsz-9"},{"Type":"text","Value":"\n"},{"Type":"link","Value":"es:Aclis-9"},{"Type":"text","Value":"\n"},{"Type":"link","Value":"Category:Materials"}]`},
		{"#REDIRECT [[The Jackal]]", `[{"Type":"text","Value":"#REDIRECT "},{"Type":"link","Value":"The Jackal"}]`}, // handle redirects gracefully
		{"[[File:Agea During Terraforming.jpg|thumb|Agea during the terraforming of Mars]]\n<p class=\"MsoNormal\">'''Agea''' is the capital city of [[Mars]]. It is located in the [[Valles Marineris]], the largest canyon in the Solar System.  </p>\n\n<p class=\"MsoNormal\">The citadel of the [[ArchGovernor]], [[Nero au Augustus]], is located in Agea, as is the [[House Bellona]] family estate. Nero holds court within the city's Forum.</p>\n\n<p class=\"MsoNormal\">Described as beautiful but strange, life in the city is both fast and cheap. It is well known for its nightlife, with rooftop nightclubs, [[gravMixers]] and [[NoiseBubbles]].</p>\n\n<p class=\"MsoNormal\">Agea is home to the [[Agea Martial Club]], where [[Gold]]'s are able to duel and practice their fighting skills against their peers.</p>\n\n<p class=\"MsoNormal\">Once [[The Institute]] trials are complete, a grand festival is held in Agea to honour the newly graduated [[Peerless Scarred]].</p>\n\n== Gallery ==\n<gallery>\nFile:Agea mars 720 pce.jpg|Agea, Mars 720 PCE\nFile:Agea Dueling Club.jpg|A dueling club in Agea, 725 PCE\nFile:725 Agea Skyline at Night.jpg|The Agea skyline at night, 725 PCE\nFile:Agea Forum.jpg|Agea's Forum, 725 PCE\n</gallery>\n[[es:Agea]]\n[[hu:Égea]]\n[[Category:Locations]]\n[[Category:Cities]]", `[{"Type":"link","Value":"File:Agea During Terraforming.jpg|thumb|Agea during the terraforming of Mars"},{"Type":"text","Value":"\nAgea is the capital city of "},{"Type":"link","Value":"Mars"},{"Type":"text","Value":". It is located in the "},{"Type":"link","Value":"Valles Marineris"},{"Type":"text","Value":", the largest canyon in the Solar System.  \n\nThe citadel of the "},{"Type":"link","Value":"ArchGovernor"},{"Type":"text","Value":", "},{"Type":"link","Value":"Nero au Augustus"},{"Type":"text","Value":", is located in Agea, as is the "},{"Type":"link","Value":"House Bellona"},{"Type":"text","Value":" family estate. Nero holds court within the city's Forum.\n\nDescribed as beautiful but strange, life in the city is both fast and cheap. It is well known for its nightlife, with rooftop nightclubs, "},{"Type":"link","Value":"gravMixers"},{"Type":"text","Value":" and "},{"Type":"link","Value":"NoiseBubbles"},{"Type":"text","Value":".\n\nAgea is home to the "},{"Type":"link","Value":"Agea Martial Club"},{"Type":"text","Value":", where "},{"Type":"link","Value":"Gold"},{"Type":"text","Value":"'s are able to duel and practice their fighting skills against their peers.\n\nOnce "},{"Type":"link","Value":"The Institute"},{"Type":"text","Value":" trials are complete, a grand festival is held in Agea to honour the newly graduated "},{"Type":"link","Value":"Peerless Scarred"},{"Type":"text","Value":".\n\n"},{"Type":"heading","Value":" Gallery "},{"Type":"text","Value":"\n\nFile:Agea mars 720 pce.jpg|Agea, Mars 720 PCE\nFile:Agea Dueling Club.jpg|A dueling club in Agea, 725 PCE\nFile:725 Agea Skyline at Night.jpg|The Agea skyline at night, 725 PCE\nFile:Agea Forum.jpg|Agea's Forum, 725 PCE\n\n"},{"Type":"link","Value":"es:Agea"},{"Type":"text","Value":"\n"},{"Type":"link","Value":"hu:Égea"},{"Type":"text","Value":"\n"},{"Type":"link","Value":"Category:Locations"},{"Type":"text","Value":"\n"},{"Type":"link","Value":"Category:Cities"}]`},
		{"[[gold]]", `[{"Type":"link","Value":"gold"}]`},
		{"{{Gold_Character\n|title1 =  [[File:Gold_Sigil.png|left|50px]] Ajax [[File:Gold_Sigil.png|right|50px]]\n|image1 =  Ajax.png \n|caption1 = \n|alias(es) = Storm Knight (as of 754 PCE)\n|gender = Male|age = ~23 (Dark Age)\n}}\n'''Ajax au Grimmus''' is the son of [[Aja au Grimmus]] and [[Atlas au Raa]].", `[{"Type":"template","Value":"Gold_Character\n|title1   File:Gold_Sigil.png|left|50px Ajax File:Gold_Sigil.png|right|50px\n|image1   Ajax.png \n|caption1  \n|alias(es)  Storm Knight (as of 754 PCE)\n|gender  Male|age  ~23 (Dark Age)\n"},{"Type":"text","Value":"\nAjax au Grimmus is the son of "},{"Type":"link","Value":"Aja au Grimmus"},{"Type":"text","Value":" and "},{"Type":"link","Value":"Atlas au Raa"},{"Type":"text","Value":"."}]`},
		{"'''Denna''' {{pron|PR|/'dɛne/}}<ref>https://youtu.be/MPEB6NAGoYk?t=1041</ref> is the primary female figure in ''[[The Name of the Wind]]''; she is arguably the main romantic interest of Kvothe, who holds an irresistible fascination with her.", `[{"Type":"text","Value":"Denna "},{"Type":"template","Value":"pron|PR|/'dɛne/"},{"Type":"text","Value":" is the primary female figure in "},{"Type":"link","Value":"The Name of the Wind"},{"Type":"text","Value":"; she is arguably the main romantic interest of Kvothe, who holds an irresistible fascination with her."}]`},
		{"{table_content}hi", `[{"Type":"table","Value":"table_content"},{"Type":"text","Value":"hi"}]`},
		{"{{template_content}}", `[{"Type":"template","Value":"template_content"}]`},
		{"{table_content}hi{table_content}", `[{"Type":"table","Value":"table_content"},{"Type":"text","Value":"hi"},{"Type":"table","Value":"table_content"}]`},
		{"{{Infobox character\n|Appearances = {{TombsOfAtuan}}\n|Mentioned = \n}}\n'''Kossil''' is a Priestess", `[{"Type":"template","Value":"Infobox character\n|Appearances  TombsOfAtuan\n|Mentioned  \n"},{"Type":"text","Value":"\nKossil is a Priestess"}]`},
		{"== Description ==\nA physical description of Aleph", `[{"Type":"heading","Value":" Description "},{"Type":"text","Value":"\nA physical description of Aleph"}]`},
		{"Hi ==== Description ==== welcome", `[{"Type":"text","Value":"Hi "},{"Type":"heading","Value":" Description "},{"Type":"text","Value":" welcome"}]`},
		{"====Description====", `[{"Type":"heading","Value":"Description"}]`},
		{"====== Description ======", `[{"Type":"heading","Value":" Description "}]`},
		{"=Description=", `[{"Type":"heading","Value":"Description"}]`},
		{"= Description =", `[{"Type":"heading","Value":" Description "}]`},
		{"==Description==", `[{"Type":"heading","Value":"Description"}]`},
		{"======= Description =======", `[{"Type":"text","Value":"="},{"Type":"heading","Value":" Description ="}]`}, // not great but this is better than destroying text
	}

	for _, c := range cases {
		tokenizer := NewTokenizer(c[0])
		tokenizer.Tokenize()

		tokensAsJson, _ := tokenizer.tokens.Stringify()

		if tokensAsJson != c[1] {
			t.Fatalf("\nInput: %s\n\nWant:\n%s\n\nRecieved:\n%s", c[0], c[1], tokensAsJson)
		}
	}
}

func TestCleanDocument(t *testing.T) {
	cases := [][]string{
		{"<p>This is a test string</p>", "This is a test string"},
		{"<img src=\"source.com\">", ""},
		{"'''Welcome''' <img src=\"source.com\">hi", "Welcome hi"},
		{"<div class=\"top-styling\">This is a <span>test</span> string</div>", "This is a test string"},
		{"'''''Welcome''''' hi", "Welcome hi"},
		{"welcome this is a link https://link.com", "welcome this is a link "},
		{"welcome this is a linkhttps://link.com", "welcome this is a link"},
		//{`<ref>should be cleaned</ref>`, ""},
		//{`<ref>should be cleaned</ref>hi<ref>should be cleaned</ref>`, "hi"},
		//{`<ref name="test">should be cleaned</ref>Hi`, "Hi"},
	}

	for _, c := range cases {
		r := cleanDocument(c[0])
		if r != c[1] {
			t.Fatalf("\nWant:\n%s\n\nRecieved:\n%s", c[1], r)
		}
	}
}

func BenchmarkTokenizer(b *testing.B) {

	for i := 0; i < b.N; i++ {
		tokenizer := NewTokenizer(sample_wikitext_lg)
		tokenizer.Tokenize()

		tokenizer = NewTokenizer(sample_wikitext_sm)
		tokenizer.Tokenize()
	}
}
