package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Config структура для хранения конфигурации
type Config struct {
	Port         string   `envconfig:"PORT" default:"8080"`
	DBConnStr    string   `envconfig:"DB_CONN_STR" required:"true"`
	KafkaBrokers []string `envconfig:"KAFKA_BROKERS" required:"true"`
	KafkaTopic   string   `envconfig:"KAFKA_TOPIC" required:"true"`
	KafkaGroupID string   `envconfig:"KAFKA_GROUP_ID" required:"true"`
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	log.Printf("Loaded configuration: %+v", config)
	return &config, nil
}
