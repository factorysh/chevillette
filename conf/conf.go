package conf

import "time"

type Conf struct {
	Source     *Source       `yaml:"source"`
	Fluentd    *Fluentd      `yaml:"fluentd"`
	Loki       *Loki         `yaml:"loki"`
	Lumber     *Lumber       `yaml:"lumber"`
	TTL        time.Duration `yaml:"ttl"`
	AuthListen string        `yaml:"auth_listen"`
	WhiteList  []string      `yaml:"whitelist"`
}

type Loki struct {
	Url      string        `yaml:"url"`
	Query    string        `yaml:"query"`
	DelayFor time.Duration `yaml:"delay_for"`
}

type Fluentd struct {
	SharedKey string `yaml:"shared_key"`
	Listen    string `yaml:"listen"`
}

type Lumber struct {
	Listen string `yaml:"listen"`
}

type Source struct {
	PrefixLines []string `yaml:"prefix_lines"`
}

func (c *Conf) Default() {
	if c.Source == nil {
		c.Source = &Source{}
	}
	if c.Fluentd == nil {
		c.Fluentd = &Fluentd{}
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
