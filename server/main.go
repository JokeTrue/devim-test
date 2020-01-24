package main

import (
	"github.com/JokeTrue/Devim-Test-Case/shared"
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
	"log"
)

var savedNumbers []int32

func main() {
	conn, err := amqp.Dial(shared.Config.AMQPConnectionURL)
	shared.HandleError(err, "Can't connect to AMQP")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	shared.HandleError(err, "Can't create a AMQP Channel")

	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("Division", true, false, false, false, nil)
	shared.HandleError(err, "Could not declare `Division` queue")

	err = amqpChannel.Qos(1, 0, false)
	shared.HandleError(err, "Could not configure QoS")

	messageChannel, err := amqpChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	shared.HandleError(err, "Could not register consumer")

	stopChan := make(chan bool)
	go func() {
		for msg := range messageChannel {
			divisionTask := &shared.DivisionTask{}

			err := proto.Unmarshal(msg.Body, divisionTask)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}

			log.Printf("Received a number: %d", divisionTask.Number)

			if err := msg.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				if divisionTask.Number%2 == 0 {
					_ = append(savedNumbers, divisionTask.Number)
					log.Printf("Number %d passed!", divisionTask.Number)

				}
			}
		}
	}()

	// Server Stop
	<-stopChan
}
