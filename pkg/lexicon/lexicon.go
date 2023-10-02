package lexicon

import (
	"fmt"
)

type Entry struct {
	Word       string
	Definition string
}

type Lexicon map[string]Entry

func New() *Lexicon {
	return &Lexicon{}
}

func (l Lexicon) Add(e Entry) {
	l[e.Word] = e
}

func (l Lexicon) Print() {
	fmt.Println("Lexicon (definition -- word)")
	fmt.Println("-------------------------------")
	i := 1
	for _, v := range l {
		fmt.Printf("%d. %s -- %s\n", i, v.Word, v.Definition)
		i++
	}
}

// func (l Lexicon) Sort
