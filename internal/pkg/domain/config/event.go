package config

type DomainEventConfig struct {
	Server          string `properties:"server"`
	Topic           string `properties:"topic"`
	GroupId         string `properties:"group-id"`
	DeadLetterTopic string `properties:"dead-letter-topic"`
}
