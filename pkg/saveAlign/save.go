package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
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
	if line != 0 {
		panic("line not found")
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

type AlignMap map[int][]int

func filterOut(arr []int, needle int) []int {
	filtered := []int{}
	for _, val := range arr {
		if val != needle {
			filtered = append(filtered, val)
		}
	}
	return filtered
}

func alignCompression(data [][]int) [][][]int {
	/*from
	[[1,2], [2,1], [3,4], [4,3], [5,5], [6,5], [7,8], [7,7]]
	to
	[[[1], [2]], [[2], [1]], [[3], [4]], [[4], [3]], [[5, 6], [5]], [[7], [7, 8]]] */
	left_align := make(AlignMap)
	right_align := make(AlignMap)

	for _, pair := range data {
		left := pair[0]
		right := pair[1]
		left_val, left_checked := left_align[left]
		right_val, right_checked := right_align[right]

		if left_checked {
			left_align[left] = append(left_val, right)
		} else {
			left_align[left] = []int{right}
		}

		if right_checked {
			right_align[right] = append(right_val, left)
		} else {
			right_align[right] = []int{left}
		}
	}

	var aligns [][][]int

	keys := make([]int, 0)
	for k := range left_align {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, left := range keys {
		val_arr := left_align[left]
		if len(val_arr) == 1 {
			left_arr, ok := right_align[val_arr[0]]
			if !ok {
				continue
			}
			if len(left_arr) == 1 {
				aligns = append(aligns, [][]int{{left}, {val_arr[0]}})
			}
			if len(left_arr) > 1 {
				aligns = append(aligns, [][]int{left_arr, val_arr})
				for _, left_value := range left_arr {
					if left_value == left {
						delete(left_align, left_value)
					}
					right_arr, ok := left_align[left_value]
					if !ok {
						continue
					}
					if len(right_arr) == 1 {
						delete(left_align, left_value)
					} else {
						left_align[left_value] = filterOut(right_arr, val_arr[0])
					}
				}
			}

		} else if len(val_arr) > 1 {
			aligns = append(aligns, [][]int{{left}, val_arr})
		}
	}
	return aligns
}

func main() {

	data :=
		[][]int{
			{1, 2}, {2, 1}, {3, 4}, {4, 3}, {5, 5}, {6, 5}, {7, 8}, {7, 7},
		}

	data_compressed := alignCompression(data)

	line := append(AlignsToBytes(data_compressed), '\n')

	saveAtLine("internal/DB/test.txt", 3, line)

	fmt.Println("done")
}
