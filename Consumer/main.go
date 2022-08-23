package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"os"
	"os/signal"
)

const brokerAddr = "127.0.0.1:9092"
const topic = "first-topic"

func main() {
	// 连接 kafka broker，创建一个 Consumer
	consumer, err := sarama.NewConsumer([]string{brokerAddr}, nil)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	// Consumer 针对 0号partition 创建一个 partitionConsumer，来消费该 partition
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	defer partitionConsumer.Close()

	// signals 是用来捕获操作系统对程序发出的中断信号的 channel
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// 用变量 consumedNum 记录接收并消费的消息数
	consumedNum := 0

consumeLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			// 打印出接收到的消息
			fmt.Println("Topic:", msg.Topic, "\tPartition:", msg.Partition, "\tOffset:", msg.Offset, "Value:", string(msg.Value))
			consumedNum++
		case <-signals:
			// 捕获到操作系统对程序发出的中断信号，关闭 producer 并结束循环
			break consumeLoop
		}
	}

	fmt.Println("Consume messages number: ", consumedNum)
}
