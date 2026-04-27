package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port   string `yaml:"port"`
		Secret string `yaml:"secret"`
	} `yaml:"server"`
	Routes []Route `yaml:"routes"`
}

type Route struct {
	Events  []string `yaml:"events"`
	Webhook string   `yaml:"webhook"`
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	return &cfg, nil
}

func (c *Config) WebhooksFor(event string) []string {
	seen := map[string]bool{}
	var urls []string
	for _, route := range c.Routes {
		for _, e := range route.Events {
			if e == "*" || e == event {
				if !seen[route.Webhook] {
					urls = append(urls, route.Webhook)
					seen[route.Webhook] = true
				}
				break
			}
		}
	}
	return urls
}
