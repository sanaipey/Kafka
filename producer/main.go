package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"time"
)

func main() {
	//This initialization makes this an actual kakfa producer app
	//Creates a kakfa producer with it's configuration settings passed in( a broker listening on 9092)
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "broker:9092"}) //The prod uses this add for initial connection
	//kafka then redirects prod to the specidied value in kafka_advertised_listeners variable, which prod/con can use.
	if err != nil {
		panic(err)
	}
	//defers the execution of this func until surrounding func returns
	defer p.Close()

	//Delivery report handler for produced messages
	// TODO:why is it a goroutine? It's a function that runs concurrently w/other funcs
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

	// TODO: Take the scanner out, can't have it deployed to kubernetes, instead hard-code
	// name of the topic
	topic := "myTopic"
	////new scanner(obtains input(os.Stdin) of types like string)
	//scanner := bufio.NewScanner(os.Stdin)
	////Scan advances scanner to the next token(which will be available through bytes or text method)
	//for scanner.Scan() {
	//	//returns the most recent token generated by a call to scan
	//	msg := scanner.Text()
	//	p.Produce(&kafka.Message{
	//		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
	//		Value:          []byte(msg),
	////		//delivery report will be delivered to the producer's object event() channel
	////	}, nil)
	//	if msg == "exit" {
	//		break
	//	}
	////}
	//
	//if scanner.Err() != nil {
	//	fmt.Printf("Failed to read user input: %v\n", err)
	//}
	//for _, word := range []string{"This", "is", "a", "slice", "containing", "messages", "produced", "to", "kafka topic"} {
	for {
		//5-second interval between messages
		time.Sleep(5 * time.Second)
		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte("word"),
		}, nil)
		if err == nil {
			fmt.Println("OK")
		} else {
			fmt.Println(err)
		}
	}
	//wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}
