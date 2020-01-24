package main

import (
	"github.com/JokeTrue/Devim-Test-Case/shared"
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"time"
)

func main() {
	conn, err := amqp.Dial(shared.Config.AMQPConnectionURL)
	shared.HandleError(err, "Can't connect to AMQP")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	shared.HandleError(err, "Can't create a AMQP Channel")

	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("Division", true, false, false, false, nil)
	shared.HandleError(err, "Could not declare `Division` queue")

	rand.Seed(time.Now().UnixNano())

	for range time.Tick(time.Duration(rand.Intn(5-1)+1) * time.Second) {
		divisionTask := shared.DivisionTask{Number: rand.Int31n(10000)}
		body, err := proto.Marshal(&divisionTask)
		if err != nil {
			shared.HandleError(err, "Error encoding proto")
		}

		log.Printf("Sending number: %d", divisionTask.Number)

		err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})

		if err != nil {
			log.Fatalf("Error publishing message: %s", err)
		}
	}
}
