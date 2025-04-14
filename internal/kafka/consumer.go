package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	kafkaconfig "ranking_video/pkg/config/kafka"
	"sync"

	"github.com/IBM/sarama"
)

// Consumer quản lý việc nhận message từ Kafka
type Consumer struct {
	consumer  sarama.ConsumerGroup
	topics    []string
	handler   EventHandler
	ready     chan bool
	ctx       context.Context
	cancelFn  context.CancelFunc
	waitGroup sync.WaitGroup
}

// EventHandler định nghĩa hàm callback xử lý khi nhận được event
type EventHandler func(Event) error

// consumerGroupHandler cài đặt interface ConsumerGroupHandler của Sarama
type consumerGroupHandler struct {
	handler EventHandler
	ready   chan bool
}

// NewConsumer tạo mới một consumer
func NewConsumer(config *kafkaconfig.KafkaConfig) (*Consumer, error) {
	client, err := sarama.NewConsumerGroup(
		config.BootstrapServers,
		config.GroupID,
		config.GetConsumerConfig(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	consumer := &Consumer{
		consumer: client,
		topics:   []string{config.TopicName},
		ready:    make(chan bool),
		ctx:      ctx,
		cancelFn: cancel,
	}

	return consumer, nil
}

// Consume bắt đầu tiêu thụ message
func (c *Consumer) Consume(handler EventHandler) error {
	c.handler = handler
	c.waitGroup.Add(1)

	go func() {
		defer c.waitGroup.Done()
		for {
			// `Consume` should be called inside an infinite loop
			if err := c.consumer.Consume(c.ctx, c.topics, &consumerGroupHandler{
				handler: c.handler,
				ready:   c.ready,
			}); err != nil {
				log.Printf("Error from consumer: %v", err)
			}

			// Check if context was cancelled, signaling that the consumer should stop
			if c.ctx.Err() != nil {
				return
			}
		}
	}()

	// Wait until the consumer has been set up
	<-c.ready
	log.Println("Consumer is ready")

	return nil
}

// Close đóng consumer
func (c *Consumer) Close() error {
	c.cancelFn()       // Trigger cancellation
	c.waitGroup.Wait() // Wait for consumption to stop
	return c.consumer.Close()
}

// Setup runs at the beginning of a new session, before ConsumeClaim
func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(h.ready)
	return nil
}

// Cleanup runs at the end of a session, once all ConsumeClaim goroutines have exited
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages()
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var event Event
		if err := json.Unmarshal(message.Value, &event); err != nil {
			log.Printf("Error parsing message: %v", err)
			session.MarkMessage(message, "")
			continue
		}

		if err := h.handler(event); err != nil {
			log.Printf("Error handling event: %v", err)
		}

		session.MarkMessage(message, "")
	}

	return nil
}

// Chuyển Event.Data sang struct cụ thể
func ParseEventData[T any](event Event) (T, error) {
	var result T

	// Serialize map[string]interface{} thành JSON
	jsonData, err := json.Marshal(event.Data)
	if err != nil {
		return result, fmt.Errorf("error serializing map data: %w", err)
	}

	// Deserialize thành struct cụ thể
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return result, fmt.Errorf("error deserializing to struct: %w", err)
	}

	return result, nil
}
