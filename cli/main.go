package main

import (
	"context"
	"os"

	"github.com/factorysh/chevillette/auth/authrequest"
	"github.com/factorysh/chevillette/conf"
	"github.com/factorysh/chevillette/input/fluentd"
	"github.com/factorysh/chevillette/log"
	"github.com/factorysh/chevillette/memory"
	"gopkg.in/yaml.v3"
)

func main() {
	cfgPath := os.Getenv("CONFIG")
	if cfgPath == "" {
		cfgPath = "/etc/chevillette.yml"
	}
	c, err := os.Open(cfgPath)
	if err != nil {
		panic(err)
	}
	var cfg conf.Conf
	err = yaml.NewDecoder(c).Decode(&cfg)
	if err != nil {
		panic(err)
	}
	cfg.Default()

	l, err := log.NewNginxLine(cfg.Source.PrefixLines...)
	if err != nil {
		panic(err)
	}
	m := memory.New(context.TODO(), cfg.TTL)

	f := fluentd.New("nginx", l.Log, m)
	ar := authrequest.New(m)
	go ar.ListenAndServe(cfg.AuthListen)
	f.ListenAndServe(cfg.Source.FluentdListen)
}
