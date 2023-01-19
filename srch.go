package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tomachalek/vertigo/v5"
)

type Triplet struct {
	Word  string
	Lemma string
	Tag   string
	Count int64
}

func (t *Triplet) mkKey() string {
	return fmt.Sprintf("%s-%s-%s", t.Word, t.Lemma, t.Tag)
}

type Extractor struct {
	counter map[string]*Triplet
}

func (tte *Extractor) ProcStruct(st *vertigo.Structure, line int, err error) error {
	return nil
}

func (tte *Extractor) ProcStructClose(st *vertigo.StructureClose, line int, err error) error {
	return nil
}

func (tte *Extractor) ProcToken(tk *vertigo.Token, line int, err error) error {
	tag := tk.Attrs[4]
	lemma := tk.Attrs[2]
	if strings.HasPrefix(tag, "N") || strings.HasPrefix(tag, "X") {
		item := Triplet{
			Word:  strings.ToLower(tk.Attrs[0]),
			Lemma: strings.ToLower(lemma),
			Tag:   tag,
		}
		k := item.mkKey()
		curr, ok := tte.counter[k]
		if ok {
			curr.Count += item.Count

		} else {
			tte.counter[k] = &item
		}
	}
	return nil
}

func (tte *Extractor) evaluate() []Triplet {
	ans := make([]Triplet, 0, int(float32(len(tte.counter))/2))
	for _, v := range tte.counter {
		if v.Count >= 10 {
			ans = append(ans, Triplet{})
		}
	}
	return ans
}

func sortAndWriteData(items []Triplet) {
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Count > items[j].Count
	})
}

func run(filePath string) []Triplet {
	parserConf := &vertigo.ParserConf{
		InputFilePath:         filePath,
		StructAttrAccumulator: "nil",
		Encoding:              "UTF-8",
		LogProgressEachNth:    100000000,
	}
	extractor := new(Extractor)
	parserErr := vertigo.ParseVerticalFile(parserConf, extractor)
	fmt.Println("ERR: ", parserErr)
	return extractor.evaluate()
}

func main() {
	items := run("")
	sortAndWriteData(items)
}
