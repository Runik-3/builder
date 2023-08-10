package dict

type Entry struct {
	Word       string
	Definition string
}

type Dictionary struct {
	Id    string
	Entry Entry
}

func New() Dictionary {
	dd := Dictionary{Id: "test", Entry: Entry{Word: "hi", Definition: "hi"}}

	return dd
}
