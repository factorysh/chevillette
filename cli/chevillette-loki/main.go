package main

import (
	"context"
	"os"
	"time"

	"github.com/factorysh/chevillette/input/loki"
)

func main() {
	u := os.Args[1]
	l, err := loki.New(u)
	if err != nil {
		panic(err)
	}
	err = l.Tail(context.TODO(), `{name="*"}`, 5*time.Second, 42, time.Now().Add(-15*time.Minute))
	if err != nil {
		panic(err)
	}
}
