package router

import (
	"context"
	"go/hioto/config"
	"go/hioto/pkg/handler/consumer"
	messagebroker "go/hioto/pkg/handler/message_broker"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type ConsumerMqtt struct {
	InstanceName, Topic string
	handlerFunc         func([]byte)
}

type ConsumerMessageBroker struct {
	consumerHandler *consumer.ConsumerHandler
	ctx             context.Context
}

func NewConsumerMessageBroker(consumerHandler *consumer.ConsumerHandler, ctx context.Context) *ConsumerMessageBroker {
	return &ConsumerMessageBroker{
		consumerHandler: consumerHandler,
		ctx:             ctx,
	}
}

func (c *ConsumerMessageBroker) StartConsumer() {
	go func() {
		for {
			ctx, cancel := context.WithCancel(context.Background())
			cancels := c.startRoutingConsumer(ctx)

			time.Sleep(1 * time.Hour)

			log.Warn("⏰ Restarting routing consumers to ensure binding is refreshed")
			for _, c := range cancels {
				c()
			}
			cancel()
		}
	}()
}

func (c *ConsumerMessageBroker) startRoutingConsumer(ctx context.Context) []context.CancelFunc {
	var cancels []context.CancelFunc

	createCtx := func() context.Context {
		c, cancel := context.WithCancel(ctx)
		cancels = append(cancels, cancel)
		return c
	}

	for _, route := range []ConsumerMqtt{
		{
			config.MQTT_INSTANCE_NAME.GetValue(),
			config.MQTT_TOPIC_AKTUATOR.GetValue(),
			c.consumerHandler.TestingConsumeAktuator,
		},
		{
			config.MQTT_INSTANCE_NAME.GetValue(),
			config.MQTT_TOPIC_SENSOR.GetValue(),
			c.consumerHandler.ControlSensorHandler,
		},
	} {
		go messagebroker.ConsumeMQTTTopic(
			createCtx(),
			route.InstanceName,
			route.Topic,
			route.handlerFunc,
		)
	}

	return cancels
}
