package kafkaconfig

import (
	"time"

	"github.com/IBM/sarama"
)

// KafkaConfig chứa các thông số cấu hình cho Kafka
type KafkaConfig struct {
	BootstrapServers []string
	GroupID          string
	TopicName        string
}

// NewKafkaConfig tạo mới một cấu hình Kafka với các giá trị mặc định
func NewKafkaConfig(brokers []string) *KafkaConfig {
	// Nếu không có broker nào được cung cấp, sử dụng giá trị mặc định
	if len(brokers) == 0 {
		brokers = []string{"localhost:9092"}
	}
	return &KafkaConfig{
		// Sử dụng "localhost:9092" khi kết nối từ máy host đến Kafka trong container
		BootstrapServers: brokers,
		GroupID:          "my-group",
		TopicName:        "new_event",
	}
}

// GetConsumerConfig trả về cấu hình cho consumer
func (c *KafkaConfig) GetConsumerConfig() *sarama.Config {
	config := sarama.NewConfig()

	// Cấu hình phiên bản Kafka
	config.Version = sarama.V2_8_0_0

	// Cấu hình consumer group
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true

	// Cấu hình timeout
	config.Net.DialTimeout = 30 * time.Second
	config.Net.ReadTimeout = 30 * time.Second
	config.Net.WriteTimeout = 30 * time.Second

	return config
}

// GetProducerConfig trả về cấu hình cho producer
func (c *KafkaConfig) GetProducerConfig() *sarama.Config {
	config := sarama.NewConfig()

	// Cấu hình phiên bản Kafka
	config.Version = sarama.V2_8_0_0

	// Cấu hình producer
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	// Chiến lược phân vùng
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	return config
}
