package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"strconv"
)

const brokerAddr = "127.0.0.1:9092"
const topic = "first-topic"

func main() {
	// 创建同步生产者
	producer, err := sarama.NewSyncProducer([]string{brokerAddr}, nil)
	if err != nil {
		panic(err)
	}

	var successNum, errorNums int
	partitionCount := map[int32]int{}

	for count := 0; count < 100; count++ {
		message := sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder("this is test: " + strconv.Itoa(count)),
		}
		// producer 向 kafka 发送消息
		partition, _, err := producer.SendMessage(&message)
		if err != nil {
			fmt.Println("Failed to send message: ", err)
			errorNums++
		} else {
			successNum++
			partitionCount[partition]++
		}
	}

	// 关闭 producer
	err = producer.Close()
	if err != nil {
		panic(err)
	}

	// 打印统计结果
	fmt.Println("Success:", successNum, "; Errors:", errorNums)
	for key, value := range partitionCount {
		fmt.Println("partition", key, ": ", value)
	}
}
