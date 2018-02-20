package solver

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/rhino1998/wordbrain/dictionary"
	"github.com/rhino1998/wordbrain/matrix"
)

//Word describes a word and the path to make it
type Word struct {
	word string
	path []matrix.Position
}

func (w Word) String() string {
	out := w.word + "["
	for _, p := range w.path {
		out += fmt.Sprintf("(%d,%d)", p.X, p.Y)
	}
	return out + "]"
}

type node struct {
	word Word
	next map[string]*node
}

func newNode(word Word) *node {
	return &node{word: word, next: make(map[string]*node)}
}

func (n *node) add(sequence ...Word) {
	if len(sequence) == 0 {
		return
	}

	ws := sequence[0].String()
	nn, ok := n.next[ws]
	if !ok {
		nn = newNode(sequence[0])
		n.next[ws] = nn
	}

	if len(sequence) > 1 {
		nn.add(sequence[1:]...)
	}
}

func (n *node) buildString(w io.Writer, level int, seq matrix.Sequence) {
	for i := 0; i < level; i++ {
		fmt.Fprint(w, "\t")
	}

	fmt.Fprintf(w, "%s\n", n.word.word)
	for _, nn := range n.next {
		nn.buildString(w, level+1, seq.Add(nn.word.path))
	}

	// if len(n.next) == 0 {
	// 	seq.Fprintln(w, level)
	// }
}

func (n *node) String() string {
	out := n.word.word + "\n"
	for _, nn := range n.next {
		out = out + strings.Replace(nn.String(), "\n", "\n\t", -1)
	}
	return out
}

//Solution describes the potential solutions to a game
type Solution struct {
	m    matrix.Matrix
	d    *dictionary.Dictionary
	root *node
}

//Solve attempts to solve a game
func Solve(m matrix.Matrix, d *dictionary.Dictionary, wordLengths []int) *Solution {
	s := &Solution{
		m:    m,
		d:    d,
		root: newNode(Word{}),
	}
	for _, n := range solve(m, d, wordLengths) {
		s.root.next[n.word.String()] = n
	}
	return s
}

func (s *Solution) String() string {
	b := bytes.NewBuffer(nil)
	fmt.Println(s.m)
	for _, n := range s.root.next {
		n.buildString(b, 0, matrix.Sequence{s.m})
	}
	return b.String()
}

//Add adds a sequence of words to the solution set
func (s *Solution) Add(sequence ...Word) {
	s.root.add(sequence...)
}

func solve(m matrix.Matrix, d *dictionary.Dictionary, wordLengths []int) []*node {
	words := make(chan Word)
	nodes := make([]*node, 0)

	wg := sync.WaitGroup{}
	for i := range m {
		for j := range m[i] {
			if m[i][j] == matrix.EmptySpace {
				continue
			}
			wg.Add(1)
			go func(i, j int) {
				findStep(m, wordLengths[0], []matrix.Position{matrix.Position{X: i, Y: j}}, d, words)
				wg.Done()
			}(i, j)
		}
	}
	go func() {
		wg.Wait()
		close(words)
	}()

	for word := range words {
		if len(wordLengths) > 1 {
			nextNodes := solve(m.Remove(word.path...), d, wordLengths[1:])
			if len(nextNodes) == 0 {
				continue
			}
			n := newNode(word)
			for _, nn := range nextNodes {
				n.next[nn.word.String()] = nn
			}
			nodes = append(nodes, n)
		} else {
			n := newNode(word)
			nodes = append(nodes, n)
		}
	}

	return nodes
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
	d *dictionary.Dictionary, c chan Word) {
	str := takePath(m, path)
	if wordLength == 1 && d.Contains(str) {
		c <- Word{word: str, path: path}
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
