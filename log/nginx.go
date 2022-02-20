package log

import (
	"github.com/factorysh/chevillette/pattern"
	iradix "github.com/hashicorp/go-immutable-radix"
)

var (
	nginxLinePattern      pattern.Matcher
	nginxLinePatternNames map[string]int
)

func init() {
	var err error
	nginxLinePattern, err = pattern.New(`<ip> - - <_> "<method> <url> <_>" <status> <_> <_> "<_>" <_>`)
	if err != nil {
		panic(err)
	}
	nginxLinePatternNames = make(map[string]int)
	for i, name := range nginxLinePattern.Names() {
		nginxLinePatternNames[name] = i
	}
}

type Nginx struct {
	tree *iradix.Node
}

func NewNginx(prefix ...string) *Nginx {
	r := iradix.New()
	for _, p := range prefix {
		r, _, _ = r.Insert([]byte(p), new(interface{}))
	}
	return &Nginx{
		tree: r.Root(),
	}

}

func (n *Nginx) Log(line []byte) (string, error) {
	m := nginxLinePattern.Matches(line)
	if len(m) == 0 { // the line doesn't match
		return "", nil
	}
	_, _, ok := n.tree.LongestPrefix(m[2])
	if ok {
		return string(m[nginxLinePatternNames["ip"]]), nil
	}
	return "", nil
}
