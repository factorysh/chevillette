package log

import (
	"github.com/factorysh/chevillette/pattern"
	iradix "github.com/hashicorp/go-immutable-radix"
)

type Nginx struct {
	tree             *iradix.Node
	linePattern      pattern.Matcher
	linePatternNames map[string]int
}

func NewNginx(prefix ...string) (*Nginx, error) {
	r := iradix.New()
	for _, p := range prefix {
		r, _, _ = r.Insert([]byte(p), new(interface{}))
	}
	n := &Nginx{
		tree: r.Root(),
	}
	err := n.SetPattern(`<ip> - - <_> "<method> <url> <_>" <status> <_> <_> "<_>" <_>`)
	return n, err
}

func (n *Nginx) SetPattern(pttrn string) error {
	var err error
	n.linePattern, err = pattern.New(pttrn)
	if err != nil {
		return err
	}
	n.linePatternNames = make(map[string]int)
	for i, name := range n.linePattern.Names() {
		n.linePatternNames[name] = i
	}
	return nil
}

func (n *Nginx) Log(line []byte) (string, error) {
	m := n.linePattern.Matches(line)
	if len(m) == 0 { // the line doesn't match
		return "", nil
	}
	_, _, ok := n.tree.LongestPrefix(m[2])
	if ok {
		return string(m[n.linePatternNames["ip"]]), nil
	}
	return "", nil
}
