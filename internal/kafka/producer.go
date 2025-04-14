package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	kafkaconfig "ranking_video/pkg/config/kafka"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

// Producer quản lý việc gửi message đến Kafka
type Producer struct {
	producer  sarama.SyncProducer
	topicName string
}

// Event đại diện cho cấu trúc sự kiện mới
type Event struct {
	ID        string                 `json:"id"`
	EventType string                 `json:"event_type"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
}

// NewProducer tạo mới một producer
func NewProducer(config *kafkaconfig.KafkaConfig) (*Producer, error) {
	producer, err := sarama.NewSyncProducer(config.BootstrapServers, config.GetProducerConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Producer{
		producer:  producer,
		topicName: config.TopicName,
	}, nil
}

// SendEvent gửi một sự kiện mới đến Kafka
func (p *Producer) SendEvent(event Event) error {
	value, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error serializing event: %w", err)
	}

	// Tạo message để gửi
	message := &sarama.ProducerMessage{
		Topic: p.topicName,
		Key:   sarama.StringEncoder(event.ID),
		Value: sarama.ByteEncoder(value),
	}

	// Gửi message
	partition, offset, err := p.producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	log.Printf("Message sent to partition %d at offset %d", partition, offset)
	return nil
}

// Close đóng producer
func (p *Producer) Close() error {
	return p.producer.Close()
}

// Chuyển struct sang Event.Data
func NewEvent(eventType string, data interface{}) (Event, error) {
	// Serialize struct thành JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return Event{}, fmt.Errorf("error serializing data: %w", err)
	}

	// Deserialize thành map[string]interface{}
	var mapData map[string]interface{}
	if err := json.Unmarshal(jsonData, &mapData); err != nil {
		return Event{}, fmt.Errorf("error deserializing to map: %w", err)
	}

	return Event{
		ID:        uuid.New().String(),
		EventType: eventType,
		Data:      mapData,
		CreatedAt: time.Now(),
	}, nil
}
