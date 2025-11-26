package config

import (
	"os"

	"github.com/joho/godotenv"
)

type EnvKey string

const (
	PORT        EnvKey = "PORT"
	DB_PATH     EnvKey = "DB_PATH"
	MAC_ADDRESS EnvKey = "MAC_ADDRESS"

	// Exchange Broker
	EXCHANGE_DIRECT EnvKey = "EXCHANGE_DIRECT"
	EXCHANGE_TOPIC  EnvKey = "EXCHANGE_TOPIC"

	// RMQ
	RMQ_URI              EnvKey = "RMQ_URI"
	RMQ_INSTANCE         EnvKey = "RMQ_IWK_INSTANCE"
	AKTUATOR_ROUTING_KEY EnvKey = "AKTUATOR_ROUTING_KEY"

	// MQTT Local
	MQTT_HOST           EnvKey = "MQTT_HOST"
	MQTT_USERNAME       EnvKey = "MQTT_USERNAME"
	MQTT_PASSWORD       EnvKey = "MQTT_PASSWORD"
	MQTT_CLIENT_ID      EnvKey = "MQTT_CLIENT_ID"
	MQTT_INSTANCE_NAME  EnvKey = "MQTT_INSTANCE_NAME"
	MQTT_TOPIC_SENSOR   EnvKey = "MQTT_TOPIC_SENSOR"
	MQTT_TOPIC_AKTUATOR EnvKey = "MQTT_TOPIC_AKTUATOR"
)

func Load() error {
	return godotenv.Load(".env")
}

func (e EnvKey) GetValue() string {
	return os.Getenv(string(e))
}
