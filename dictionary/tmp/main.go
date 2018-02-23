package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var header = `
package dictionary
`
var wordHeader = `
var wordlist = map[string]struct{}{
`

var prefixHeader = `
var prefixlist = map[string]struct{}{
`

var footer = `
}
`

func main() {
	words, err := loadWords(os.Args[1])

	if err != nil {
		log.Println(err)
	}

	wordlist := make(map[string]struct{})
	prefixlist := make(map[string]struct{})

	for _, word := range words {
		wordlist[word] = struct{}{}
		for i := range word {
			prefixlist[word[:i]] = struct{}{}
		}
	}

	fmt.Println(header)
	fmt.Println(wordHeader)

	for word := range wordlist {
		fmt.Printf("%q: struct{}{},\n", word)
	}
	fmt.Println(footer)

	fmt.Println(prefixHeader)

	for word := range prefixlist {
		fmt.Printf("%q: struct{}{},\n", word)
	}
	fmt.Println(footer)
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
