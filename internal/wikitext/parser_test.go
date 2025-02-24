package wikitext

import (
	_ "embed"
	"strconv"
	"testing"
)

func TestParser(t *testing.T) {
	cases := [][]string{
		{"'''Denna''' {{pron|PR|/'dɛne/}}<ref>https://youtu.be/MPEB6NAGoYk?t=1041</ref> is the primary female figure in ''[[The Name of the Wind]]''; she is arguably the main romantic interest of Kvothe, who holds an irresistible fascination with her.", "Denna is the primary female figure in The Name of the Wind; she is arguably the main romantic interest of Kvothe, who holds an irresistible fascination with her."},
		{"'''Akari au Raa''' was a [[Gold]], the progenitor of [[House Raa]], and one of the founders of the [[The Society|Society]].", "Akari au Raa was a Gold, the progenitor of House Raa, and one of the founders of the Society."},
		{"#REDIRECT [[The Jackal]]", ""},
		{"[[File:Agea During Terraforming.jpg|thumb|Agea during the terraforming of Mars]]\n<p class=\"MsoNormal\">'''Agea''' is the capital city of [[Mars]]. It is located in the [[Valles Marineris]], the largest canyon in the Solar System.  </p>\n\n<p class=\"MsoNormal\">The citadel of the [[ArchGovernor]], [[Nero au Augustus]], is located in Agea, as is the [[House Bellona]] family estate. Nero holds court within the city's Forum.</p>\n\n<p class=\"MsoNormal\">Described as beautiful but strange, life in the city is both fast and cheap. It is well known for its nightlife, with rooftop nightclubs, [[gravMixers]] and [[NoiseBubbles]].</p>\n\n<p class=\"MsoNormal\">Agea is home to the [[Agea Martial Club]], where [[Gold]]'s are able to duel and practice their fighting skills against their peers.</p>\n\n<p class=\"MsoNormal\">Once [[The Institute]] trials are complete, a grand festival is held in Agea to honour the newly graduated [[Peerless Scarred]].</p>\n\n== Gallery ==\n<gallery>\nFile:Agea mars 720 pce.jpg|Agea, Mars 720 PCE\nFile:Agea Dueling Club.jpg|A dueling club in Agea, 725 PCE\nFile:725 Agea Skyline at Night.jpg|The Agea skyline at night, 725 PCE\nFile:Agea Forum.jpg|Agea's Forum, 725 PCE\n</gallery>\n[[es:Agea]]\n[[hu:Égea]]\n[[Category:Locations]]\n[[Category:Cities]]", "Agea is the capital city of Mars."},
		{"{{Gold_Character\n|title1 =  [[File:Gold_Sigil.png|left|50px]] Ajax [[File:Gold_Sigil.png|right|50px]]\n|image1 =  Ajax.png \n|caption1 = \n|alias(es) = Storm Knight (as of 754 PCE)\n|gender = Male|age = ~23 (Dark Age)\n}}\n'''Ajax au Grimmus''' is the son of [[Aja au Grimmus]] and [[Atlas au Raa]].", "Ajax au Grimmus is the son of Aja au Grimmus and Atlas au Raa."},
		{"[[File:Inn.jpg|thumb]][[Lamia]] created the '''Inn''' to trick [[Yvaine]] into believing that she was a simple Inn keeper's Wife.", "Lamia created the Inn to trick Yvaine into believing that she was a simple Inn keeper's Wife."},
		{"{{quote|We found him wandering around, with a candle.|A scriv about [[Kvothe]]{{ref|TNOTW}}}}\n\n[[File:Kvothe by ladyroadx-d4hvaki.jpg|thumb|290px|Fan Art by\nhttp://ladyroadx.deviantart.com/]]\n\nA '''Scriv''' is a student who works under [[Master Lorren]], specifically in [[University|The University's]] [[The Archives|Archives]].", "A Scriv is a student who works under Master Lorren, specifically in The University's Archives."},
		{"{{Infobox character\n|Appearances = {{TombsOfAtuan}}\n|Mentioned = \n}}\n'''Kossil''' is a Priestess", "Kossil is a Priestess"},
		{"A.203 is a robot from the achovy system", "A.203 is a robot from the achovy system"},
		{"==Heading==A.203 is a robot from the achovy system", "A.203 is a robot from the achovy system"},
		{"handwritten in [[wikipedia:Octavo|octavo-sized pages]]", "handwritten in octavo-sized pages"},
		{"[[Category:Poems]]", ""},
		{"[[Media:Poems]]", ""},
		{"[[Media:Poems|Display should not show]]", ""},
		{"[[Category:Poems|Display should show]]", "Display should show"},
		{`{{quote|Simmon [...] had a great deal to offer. He was a gemstone in the rough. Not stunning at first glance, but with a great deal of worth beneath the surface. Sim was tender, kind, and attentive [...] Sim was a prince.|[[Kvothe]]{{ref|TWMF}}}}
{{character infobox
|image = Playing Cards poster Simmon.png
|alias = Sim
|family = *Duke of Dalonir (father)
*Five unnamed siblings
|species = Human
|gender = Male
|hair = Sandy
|eye = Light
|ethnicity = [[The Aturan Empire|Aturan]]
|occupation = [[University]] student
|field = Alchemy, Poetry
|rank = Re'lar
|skin = Pale}}

'''Simmon''' {{pron|PR|/'sɪmən/}}<ref>https://youtu.be/MPEB6NAGoYk?t=101</ref>, usually referred to as Sim, is one of [[Kvothe|Kvothe's]] best friends at [[the University]].`,
			"Simmon , usually referred to as Sim, is one of Kvothe's best friends at the University.",
		},
	}

	for _, c := range cases {
		r, _ := ParseDefinition(c[0], 1)
		if r != c[1] {
			t.Fatalf("\nWant:\n%s\n\nRecieved:\n%s", c[1], r)
		}
	}
}

func TestParserDepth(t *testing.T) {
	cases := [][]string{
		{"'''Kossil''' is a Priestess at the Place of the Tombs of Atuan. She serves as the High Priestess of the Godking. She, along with [[Thar]] and [[Tenar]], are the high authorities there.", "Kossil is a Priestess at the Place of the Tombs of Atuan.", "1"},
		{"'''Kossil''' is a Priestess at the Place of the Tombs of Atuan. She serves as the High Priestess of the Godking. She, along with [[Thar]] and [[Tenar]], are the high authorities there.", "Kossil is a Priestess at the Place of the Tombs of Atuan. She serves as the High Priestess of the Godking.", "2"},
		{"'''Kossil''' is a Priestess at the Place of the Tombs of Atuan. She serves as the High Priestess of the Godking. She, along with [[Thar]] and [[Tenar]], are the high authorities there.", "Kossil is a Priestess at the Place of the Tombs of Atuan. She serves as the High Priestess of the Godking. She, along with Thar and Tenar, are the high authorities there.", "3"},
	}

	for _, c := range cases {
		depth, _ := strconv.Atoi(c[2])
		r, _ := ParseDefinition(c[0], depth)
		if r != c[1] {
			t.Fatalf("\nWant:\n%s\n\nRecieved:\n%s", c[1], r)
		}
	}
}

func BenchmarkParsing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseDefinition(sample_wikitext_lg, 1)
		ParseDefinition(sample_wikitext_sm, 1)
	}
}
