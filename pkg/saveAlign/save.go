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
		if line == 0 {
			offset = i
		}

		if b == '\n' {
			line--
			if line == -1 {
				next_line = i
				break
			}
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

	f, err := os.OpenFile("internal/DB/test.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	end_line := offset + len(data)

	f.WriteAt(data, int64(offset))
	f.WriteAt(content[next_line:], int64(end_line))

}

func main() {
	f, err := os.Create("internal/DB/test.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data :=
		[][][][]int{
			{{{1, 2}, {2}}, {{3}, {1, 4}}, {{5, 4}, {6}}, {{7, 6}, {8}}},
			{{{0}, {0}}, {{1}, {1}}, {{2}, {2}}, {{3}, {3}}, {{4}, {4}}, {{5}, {5}}, {{6}, {6}}},
			{{{0}, {1}}, {{1}, {2}}, {{2}, {3}}, {{3}, {4}}, {{4}, {5}}, {{5}, {6}}, {{6}, {7}}},
		}

	for i := 0; i < len(data); i++ {
		text := AlignsToBytes(data[i])
		_, err = f.WriteString(string(text) + "\n")
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("done")
}
