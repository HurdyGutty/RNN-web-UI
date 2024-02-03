package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func AlignsToBytes(align [][][]int) []byte {
	text, err := json.Marshal(align)
	if err != nil {
		panic(err)
	}
	return text
}

func offsetLines(data []byte, line int) (int, int) {
	offset := 0
	next_line := 0
	for i, b := range data {

		if b == '\n' {

			line--
			if line == 0 {
				offset = i + 1
			}

			if line == -1 {
				next_line = i + 1
				break
			}
			fmt.Println(line, offset, next_line)
		}
	}
	return offset, next_line
}

func saveAtLine(path string, line int, data []byte) {
	content, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	offset, next_line := offsetLines(content, line)

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fmt.Println(offset, next_line)

	end_line := offset + len(data)
	f.WriteAt(content[:offset], 0)
	f.WriteAt(data, int64(offset))
	if next_line != 0 {
		f.WriteAt(content[next_line:], int64(end_line))
	} else {
		f.Truncate(int64(end_line - 1))
	}

}

func main() {

	data :=
		[][][]int{
			{{1, 2}, {2}}, {{3}, {1, 4}}, {{5, 4}, {6}}, {{7, 6}, {8}},
		}

	line := append(AlignsToBytes(data), '\n')

	saveAtLine("internal/DB/test.txt", 0, line)

	fmt.Println("done")
}
