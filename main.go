package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/rhino1998/wordbrain/dictionary"
	"github.com/rhino1998/wordbrain/matrix"
	"github.com/rhino1998/wordbrain/solver"
)

func main() {
	m := matrix.Matrix{
		[]matrix.Space{'p', 'e', 'e', 'r'},
		[]matrix.Space{'a', 'h', 'd', 'b'},
		[]matrix.Space{'e', 't', 'c', 'a'},
		[]matrix.Space{'m', 's', 't', 'a'},
	}
	words, err := loadWords("words.txt")
	if err != nil {
		log.Fatal(err)
	}
	d := dictionary.New(words...)
	fmt.Println(m)
	s := solver.Solve(m, d, []int{4, 5, 3, 4})
	fmt.Println(s)
}

func loadWords(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
