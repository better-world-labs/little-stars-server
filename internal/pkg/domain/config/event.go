package config

type DomainEventConfig struct {
	Server          string `yaml:"server"`
	Topic           string `yaml:"topic"`
	GroupId         string `yaml:"group-id"`
	DeadLetterTopic string `yaml:"dead-letter-topic"`
}
