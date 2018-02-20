package matrix

import (
	"bytes"
	"fmt"
	"io"
)

type Sequence []Matrix

func (s Sequence) Add(path []Position) Sequence {
	seq := append(make(Sequence, 0, len(s)+1), s...)
	return append(seq, seq[len(s)-1].Remove(path...))
}

func (s Sequence) Fprintln(w io.Writer, level int) {
	for _, m := range s {
		for i := 0; i < level; i++ {
			fmt.Fprintf(w, "\t")
		}
		for i := range m {
			fmt.Fprintf(w, "| ")
			for char := range m[len(m)-i-1] {
				fmt.Fprintf(w, "%s ", string(char))
			}
			fmt.Fprintf(w, "|")
		}
		fmt.Fprintf(w, "\n")
	}
}

func (s Sequence) String() string {
	b := bytes.NewBuffer(nil)
	s.Fprintln(b, 0)
	return b.String()
}
