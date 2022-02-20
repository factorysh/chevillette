package conf

import "time"

type Conf struct {
	Source     Source        `yaml:"source"`
	TTL        time.Duration `yaml:"ttl"`
	AuthListen string        `yaml:"auth_listen"`
}

type Source struct {
	PrefixLines   []string `yaml:"prefix_lines"`
	FluentdListen string   `yaml:"fluentd_listen"`
}
