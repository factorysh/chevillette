package main

import (
	"context"
	"os"
	"time"

	"github.com/factorysh/chevillette/auth/authrequest"
	"github.com/factorysh/chevillette/conf"
	"github.com/factorysh/chevillette/input/fluentd"
	"github.com/factorysh/chevillette/input/loki"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m := memory.New(ctx, cfg.TTL)

	ar := authrequest.New(m)

	if cfg.Fluentd != nil && cfg.Fluentd.Listen != "" {
		f, err := fluentd.New("nginx", l.Log, m, cfg.Fluentd)
		if err != nil {
			panic(err)
		}
		go f.ListenAndServe(cfg.Fluentd.Listen)
	}

	if cfg.Loki != nil {
		l, err := loki.New(cfg.Loki.Url)
		if err != nil {
			panic(err)
		}
		go func() {
			err = l.Tail(context.TODO(), cfg.Loki.Query, 10*time.Second, 1000, time.Now())
			if err != nil {
				panic(err)
			}
		}()
	}

	err = ar.ListenAndServe(cfg.AuthListen)
	if err != nil {
		panic(err)
	}
}
