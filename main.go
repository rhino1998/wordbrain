package main

import (
	"bufio"
	"fmt"
	"math"
	"os"

	"github.com/rhino1998/wordbrain/dictionary"
	"github.com/rhino1998/wordbrain/matrix"
	"github.com/rhino1998/wordbrain/solver"
)

func main() {

	m := makeMatrix(os.Args[1])
	d := dictionary.NewDefault()
	fmt.Println(m)
	s := solver.Solve(m, d, getWordLengths(os.Args[2]))
	fmt.Println(s)
}

func makeMatrix(strMatrix string) matrix.Matrix {
	dims := int(math.Sqrt(float64(len(strMatrix))))
	m := make(matrix.Matrix, dims)
	for i := 0; i < dims; i++ {
		m[i] = make([]matrix.Space, dims)
		for j := 0; j < dims; j++ {
			m[i][j] = matrix.Space(strMatrix[(len(strMatrix))-(i+1)*dims+j])
		}
	}
	return m
}

func getWordLengths(strLengths string) []int {
	lengths := make([]int, 0, 4)
	for _, length := range strLengths {
		lengths = append(lengths, int(length-'0'))
	}
	return lengths
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
