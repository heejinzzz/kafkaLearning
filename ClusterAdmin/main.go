package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

const brokerAddr = "127.0.0.1:9092"

func main() {
	// 连接 kafka 集群并创建 ClusterAdmin
	admin, err := sarama.NewClusterAdmin([]string{brokerAddr}, nil)
	if err != nil {
		panic(err)
	}

	// 列出所有 topic
	topics, err := admin.ListTopics()
	if err != nil {
		panic(err)
	}
	fmt.Println(len(topics), " topics in cluster:")
	for key := range topics {
		fmt.Println("Topic: ", key)
	}

	// 创建 topic
	topicDetail := sarama.TopicDetail{
		NumPartitions:     3,
		ReplicationFactor: 2,
	}
	err = admin.CreateTopic("test-topic", &topicDetail, false)
	if err != nil {
		panic(err)
	}

	// 查看指定的 topic 的详细信息
	topicMetadata, err := admin.DescribeTopics([]string{"test-topic"})
	if err != nil {
		panic(err)
	}
	fmt.Println("TopicName: ", topicMetadata[0].Name)
	fmt.Println("PartitionCount: ", len(topicMetadata[0].Partitions))
	fmt.Println("ReplicationFactor: ", len(topicMetadata[0].Partitions[0].Replicas))
	fmt.Println("ReplicaAssignment: ")
	for _, partitionMetadata := range topicMetadata[0].Partitions {
		fmt.Println("\tpartition:", partitionMetadata.ID, "\tleader:", partitionMetadata.Leader, "\treplicas:", partitionMetadata.Replicas, "\tisr:", partitionMetadata.Isr)
	}

	// 删除指定的 topic
	err = admin.DeleteTopic("test-topic")
	if err != nil {
		panic(err)
	}
}
