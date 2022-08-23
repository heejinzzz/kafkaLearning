package main

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
)

const brokerAddr = "127.0.0.1:9092"
const topic = "first-topic"
const groupId = "test_consumerGroup"

// 注册 ConsumerGroupHandler 结构体并实现 interface 中定义的方法
type consumerGroupHandler struct{}

func (handler consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (handler consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (handler consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Println("Topic:", msg.Topic, "\tPartition:", msg.Partition, "\tOffset:", msg.Offset, "Value:", string(msg.Value))
		session.MarkMessage(msg, "") // MarkMessage 会将一个 message 标记为已消费
	}
	return nil
}

func main() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// 连接 kafka broker，创建一个 ConsumerGroup
	consumerGroup, err := sarama.NewConsumerGroup([]string{brokerAddr}, groupId, nil)
	if err != nil {
		panic(err)
	}
	defer consumerGroup.Close()

	// 开启一个协程来记录 consumerGroup 接收消息消费过程中出现的错误
	go func() {
		for err := range consumerGroup.Errors() {
			fmt.Println(err)
		}
	}()

	// 创建一个 ConsumerGroupHandler
	handler := consumerGroupHandler{}

	err = consumerGroup.Consume(context.Background(), []string{topic}, handler)
	if err != nil {
		panic(err)
	}
}
