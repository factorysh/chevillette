package conf

import "time"

type Conf struct {
	Source     *Source       `yaml:"source"`
	TTL        time.Duration `yaml:"ttl"`
	AuthListen string        `yaml:"auth_listen"`
	WhiteList  []string      `yaml:"whitelist"`
}

type Source struct {
	PrefixLines   []string `yaml:"prefix_lines"`
	FluentdListen string   `yaml:"fluentd_listen"`
}

func (c *Conf) Default() {
	if c.Source == nil {
		c.Source = &Source{}
	}
	if c.Source.FluentdListen == "" {
		c.Source.FluentdListen = "127.0.0.1:24224"
	}
	if c.Source.PrefixLines == nil {
		c.Source.PrefixLines = []string{}
	}
	if c.TTL == 0 {
		c.TTL = time.Hour
	}
	if c.AuthListen == "" {
		c.AuthListen = "127.0.0.1:8080"
	}
	if c.WhiteList == nil {
		c.WhiteList = []string{}
	}
}
