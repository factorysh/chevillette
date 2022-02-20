package main

import (
	"context"
	"time"

	"github.com/factorysh/chevillette/auth/authrequest"
	"github.com/factorysh/chevillette/input/fluentd"
	"github.com/factorysh/chevillette/log"
	"github.com/factorysh/chevillette/memory"
)

func main() {
	l, err := log.NewNginxLine("/")
	if err != nil {
		panic(err)
	}
	m := memory.New(context.TODO(), time.Hour)

	f := fluentd.New("nginx", l.Log, m)
	ar := authrequest.New(m)
	go ar.ListenAndServe("0.0.0.0:8080")
	f.ListenAndServe("0.0.0.0:24224")
}
