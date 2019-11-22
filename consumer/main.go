package main

import (
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func main() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest", //action to take when there's no initial offset(like a bookmark)
		// It's how consumers read messages independently,resets to smallest offset
	})
	if err != nil {
		panic(err)
	}
	// Subscribes to the provided list of topics
	// nil = no rebalance event call back
	// regexpression anything that includes topic
	err = c.SubscribeTopics([]string{"myTopic", "^aRegex.*[Tt]opic"}, nil)

	for {
		//ReadMessage polls the consumer for a message, timeout` may be set to -1 for indefinite wait.
		msg, err := c.ReadMessage(-1)
		if err == nil {
			//msg.TopicPartition provides partition-specific information (such as topic, partition and offset)
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			if string(msg.Value) == "exit" {
				break
			} else {
				// The client will automatically try to recover from all errors.
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}
	}
	//Close consumer instance
	c.Close()
}
