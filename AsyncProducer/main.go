package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
)

const brokerAddr = "127.0.0.1:9092"
const topic = "first-topic"

func main() {
	// 设置要创建的生产者的参数
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true                // 要求 Producer 返回消息发送成功的信息
	config.Producer.Compression = sarama.CompressionSnappy // 设置数据压缩格式为 snappy
	config.Producer.RequiredAcks = sarama.WaitForAll       // 设定当 leader 和 ISR 队列均收到消息时才发送确认帧，防止数据丢失

	// 连接 kafka broker ，创建异步生产者
	producer, err := sarama.NewAsyncProducer([]string{brokerAddr}, config)
	if err != nil {
		panic(err)
	}

	// signals 是用来捕获操作系统对程序发出的中断信号的 channel
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var wg sync.WaitGroup
	var enqueueNum, errorNum, successNum int
	partitionCount := map[int32]int{}

	// 开启一个协程来记录 producer 向 kafka 发送消息成功次数，以及消息分发到了哪个 partition 上
	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range producer.Successes() {
			successNum++
			partitionCount[msg.Partition]++
		}
	}()

	// 开启一个协程来读取 producer 向 kafka 发送消息过程中出现的错误
	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range producer.Errors() {
			fmt.Println(err)
			errorNum++
		}
	}()

	// producer 向 kafka 发送消息
ProducerLoop:
	for count := 0; ; count++ {
		time.Sleep(10 * time.Millisecond)
		message := sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder("this is test: " + strconv.Itoa(count)),
		}
		select {
		case producer.Input() <- &message:
			enqueueNum++
		case <-signals:
			// 捕获到操作系统对程序发出的中断信号，关闭 producer 并结束循环
			producer.AsyncClose()
			break ProducerLoop
		}
	}

	wg.Wait()

	// 打印统计结果
	fmt.Println("Enqueued:", enqueueNum, "; Success:", successNum, "; Errors:", errorNum)
	for key, value := range partitionCount {
		fmt.Println("partition", key, ": ", value)
	}
}
