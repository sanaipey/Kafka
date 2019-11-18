package main

import (
	"bufio"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

func main() {
	//This initialization makes this an actual kakfa producer app
	//Creates a kakfa producer with it's configuration settings passed in( a broker listening on 29092)
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		panic(err)
	}
	//defers the execution of this func until surrounding func returns
	defer p.Close()

	//Delivery report handler for produced messages
	// TODO:why is it a goroutine?
	go func() {
		for e := range p.Events() {
			//type assertion(to find the type of that one event(e)in the event's stream: if e is of type Message( the if statements will be executed))
			// this switch stmt returns the first case equal to the condition expression
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivered failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// name of the topic
	topic := "myTopic"
	//new scanner(obtains input(os.Stdin) of types like string)
	scanner := bufio.NewScanner(os.Stdin)
	//Scan advances scanner to the next token(which will be available through bytes or text method)
	for scanner.Scan() {
		//returns the most recent token generated by a call to scan
		msg := scanner.Text()
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(msg),
			//delivery report will be delivered to the producer's object event() channel
		}, nil)
		if msg == "exit" {
			break
		}
	}

	if scanner.Err() != nil {
		fmt.Printf("Failed to read user input: %v\n", err)
	}

	//wait for message deliveries before shutting down
	p.Flush(15 * 1000)

}
