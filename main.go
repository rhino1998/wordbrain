package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/rhino1998/wordbrain/dictionary"
	"github.com/rhino1998/wordbrain/matrix"
)

func main() {
	m := matrix.Matrix{
		[]matrix.Space{'n', 'c', 'b', 'u'},
		[]matrix.Space{'o', 'b', 't', 't'},
		[]matrix.Space{'r', 'i', 'a', 'e'},
		[]matrix.Space{'d', 'k', 'n', 'r'},
	}
	words, err := loadWords("words.txt")
	if err != nil {
		log.Fatal(err)
	}
	d := dictionary.New(words...)
	fmt.Println(m)
	fmt.Println(solve(m, d, []int{6, 5, 5}, 0))
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

func solve(m matrix.Matrix, d *dictionary.Dictionary, wordLengths []int, level int) (string, int) {
	paths := make(chan []matrix.Position)

	out := ""

	wg := sync.WaitGroup{}
	for i := range m {
		for j := range m[i] {
			wg.Add(1)
			go func(i, j int) {
				findStep(m, wordLengths[0], []matrix.Position{matrix.Position{i, j}}, d, paths)
				wg.Done()
			}(i, j)
		}
	}
	go func() {
		wg.Wait()
		close(paths)
	}()

	numResults := 0
	for path := range paths {
		tmpOut := ""
		numResults++
		for i := 0; i < level; i++ {
			tmpOut += "\t"
		}
		if len(wordLengths) == 1 {
			tmpOut += fmt.Sprintf("%s\n", takePath(m, path))
		} else {
			strResults, nestResults := solve(m.Remove(path...), d, wordLengths[1:], level+1)
			if nestResults == 0 {
				continue
			}
			tmpOut += fmt.Sprintf(
				"%s\n%s\n%s",
				takePath(m, path),
				m.Remove(path...),
				strResults,
			)
		}
		out += tmpOut
	}
	return out, numResults
}

func takePath(m matrix.Matrix, path []matrix.Position) string {
	word := ""
	for _, pos := range path {
		word += string(m[pos.X][pos.Y])
	}
	return word
}

func posInPath(p matrix.Position, path []matrix.Position) bool {
	for _, op := range path {
		if p == op {
			return true
		}
	}
	return false
}

func findStep(m matrix.Matrix, wordLength int, path []matrix.Position,
	d *dictionary.Dictionary, c chan []matrix.Position) {
	if wordLength == 1 && d.Contains(takePath(m, path)) {
		c <- path
		return
	}
	wg := sync.WaitGroup{}
	for _, pos := range m.ValidMoves(path[len(path)-1]) {

		if posInPath(pos, path) {
			continue
		}

		np := append(make([]matrix.Position, 0, len(path)+1), path...)
		np = append(np, pos)

		if !d.ContainsPrefix(takePath(m, np)) {
			continue
		}

		wg.Add(1)
		go func() {
			findStep(m, wordLength-1, np, d, c)
			wg.Done()
		}()
	}
	wg.Wait()
}
