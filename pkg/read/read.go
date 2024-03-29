package read

import (
	"os"
)

func readData(path string) []byte {
	content, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	return content
}

func parseLanguage(data []byte) [][]string {
	var lines [][]string
	var line []string
	var word []byte

	appendWord := func(word []byte) []byte {
		if string(word) != "" {
			line = append(line, string(word))
		}
		return []byte{}
	}

	for _, b := range data {

		if b == '\r' {
			continue
		}
		if b == ',' || b == '.' || b == '?' || b == '!' || b == ':' || b == ';' {
			continue
		}
		if b == '\n' {
			word = appendWord(word)
			lines = append(lines, line)
			line = []string{}
			continue
		}

		if b == ' ' {
			word = appendWord(word)
			continue
		}

		word = append(word, b)

	}
	appendWord(word)
	lines = append(lines, line)
	return lines
}

type Dict map[string]interface{}
type AlignmentValues Dict

func newDict() AlignmentValues {
	return AlignmentValues{
		"Nom":   Dict{"Key": "Nom", "Data": []string{}},
		"Eng":   Dict{"Key": "Eng", "Data": []string{}},
		"Align": [][]int{},
	}
}

func (values AlignmentValues) mockData(nom, eng []string, align [][]int) AlignmentValues {
	values["Nom"] = Dict{"Key": "Nom", "Data": nom}
	values["Eng"] = Dict{"Key": "Eng", "Data": eng}
	values["Align"] = align
	return values
}

type Page struct {
	Page      int
	TotalPage int
	Values    AlignmentValues
}

func newPage(page, totalPage int, values AlignmentValues) Page {
	return Page{
		Page:      page,
		TotalPage: totalPage,
		Values:    values,
	}
}

type Pages []Page

func newPages() Pages {
	newPages := []Page{}
	return newPages
}

func MockPages() Pages {
	newPages := []Page{
		newPage(1, 3, newDict().mockData(
			[]string{"Je", "parle", "francais"},
			[]string{"I", "speak", "French"},
			[][]int{{0, 0}, {1, 1}, {2, 2}})),
		newPage(2, 3, newDict().mockData(
			[]string{"Battre", "le", "fer", "pendant", "qu'il", "est", "chaud"},
			[]string{"Strike", "the", "iron", "while", "it", "is", "hot"},
			[][]int{{0, 0}, {1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}, {6, 6}})),
		newPage(3, 3, newDict().mockData(
			[]string{"En", "faire", "tout", "un", "fromage"},
			[]string{"To", "make", "a", "whole", "cheese"},
			[][]int{{0, 0}, {1, 1}, {2, 3}, {3, 2}, {4, 4}})),
	}
	return newPages
}

func parseAlign(data []byte) [][][]int {
	var alignments [][][]int
	var alignment_line [][]int
	var isRight bool = false
	var left []int
	var right []int
	var level int = 0

	for _, b := range data {
		if b == '\n' {
			alignments = append(alignments, alignment_line)
			alignment_line = [][]int{}
			left = []int{}
			right = []int{}
			level = 0
			isRight = false
			continue
		}
		if b == '\r' {
			continue
		}

		if b == ',' {
			switch level {
			case 1:
				for _, v := range left {
					for _, w := range right {
						alignment_line = append(alignment_line, []int{v, w})
					}
				}
				left = []int{}
				right = []int{}
				isRight = false
			case 2:
				isRight = !isRight
			}
			continue
		}

		if b == '[' {
			level += 1
			continue
		}

		if b == ']' {
			level -= 1
			if level == 0 {
				for _, v := range left {
					for _, w := range right {
						alignment_line = append(alignment_line, []int{v, w})
					}
				}
				left = []int{}
				right = []int{}
				isRight = false
			}
			continue
		}

		if (b >= '0' && b <= '9') && level == 3 {
			if isRight {
				right = append(right, int(b-'0')-1)
			} else {
				left = append(left, int(b-'0')-1)
			}
		}
	}
	alignments = append(alignments, alignment_line)
	return alignments
}

func ParseData() Pages {
	var static_path string = "internal/DB/"
	nom_data := readData(static_path + "test-vie.txt")
	eng_data := readData(static_path + "test-eng.txt")
	align_data := readData(static_path + "aligned_vie-eng.txt")

	nom := parseLanguage(nom_data)
	eng := parseLanguage(eng_data)
	align := parseAlign(align_data)
	pages := newPages()
	length := len(nom)

	for i := 0; i < length; i++ {
		pages = append(pages, newPage(i+1, length, newDict().mockData(nom[i], eng[i], align[i])))
	}

	return pages
}

// func main() {
// 	fmt.Printf("%v\n", ParseData())
// }
