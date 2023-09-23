package dict

import (
	"fmt"
)

type Entry struct {
	Word       string
	Definition string
}

type Dict map[string]Entry

func New() *Dict {
	dict := Dict{}

	return &dict
}

func (d Dict) Add(e Entry) {
	d[e.Word] = e
}

func (d Dict) Print() {
	fmt.Println("Dictionary (definition -- word)")
	fmt.Println("-------------------------------")
	i := 1
	for _, v := range d {
		fmt.Printf("%d. %s -- %s\n", i, v.Word, "")
		i++
	}
}
