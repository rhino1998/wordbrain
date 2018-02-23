package dictionary

//Dictionary contains valid words
type Dictionary struct {
	valid       map[string]struct{}
	prefixValid map[string]struct{}
}

//New returns a new dictionary with words
func New(words ...string) *Dictionary {
	d := &Dictionary{prefixValid: make(map[string]struct{}), valid: make(map[string]struct{})}
	for _, word := range words {
		d.valid[word] = struct{}{}
		wordRune := []rune(word)
		for i := range wordRune {
			d.prefixValid[string(wordRune[:i])] = struct{}{}
		}
	}
	return d
}

func NewDefault() *Dictionary {
	return &Dictionary{valid:wordlist, prefixValid:prefixlist}
}

//Contains checks for the presence of word in d
func (d *Dictionary) Contains(word string) bool {
	_, ok := d.valid[word]
	return ok
}
func (d *Dictionary) ContainsPrefix(word string) bool {
	_, ok := d.prefixValid[word]
	return ok
}

type node struct {
	suffixes map[rune]*node
}

func newNode() *node {
	return &node{suffixes: make(map[rune]*node)}
}

func (n *node) add(suffix []rune) {
	if len(suffix) == 0 {
		return
	}

	if nn, ok := n.suffixes[suffix[0]]; ok {
		nn.add(suffix[1:])
		return
	}
	nn := newNode()
	n.suffixes[suffix[0]] = nn
	nn.add(suffix[1:])
}

func (n *node) contains(suffix []rune) bool {
	if len(suffix) == 0 {
		return true
	}
	if nn, ok := n.suffixes[suffix[0]]; ok {
		return nn.contains(suffix[1:])
	}
	return false
}
