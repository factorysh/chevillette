package main

import (
	"github.com/factorysh/chevillette/input/fluentd"
	"github.com/factorysh/chevillette/log"
)

func main() {
	l, err := log.NewNginxLine("/")
	if err != nil {
		panic(err)
	}
	f := fluentd.New("nginx", l.Log)
	f.ListenAndServe("0.0.0.0:24224")
}
