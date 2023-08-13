package dict

import "fmt"

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
	fmt.Println(d)
}
