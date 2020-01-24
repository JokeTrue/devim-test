package main

import (
	devim_case "devim-case"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func main() {
	conn, err := amqp.Dial(devim_case.Config.AMQPConnectionURL)
	devim_case.HandleError(err, "Can't connect to AMQP")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	devim_case.HandleError(err, "Can't create a AMQP Channel")

	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("Division", true, false, false, false, nil)
	devim_case.HandleError(err, "Could not declare `Division` queue")

	err = amqpChannel.Qos(1, 0, false)
	devim_case.HandleError(err, "Could not configure QoS")

	messageChannel, err := amqpChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	devim_case.HandleError(err, "Could not register consumer")

	stopChan := make(chan bool)

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range messageChannel {
			log.Printf("Received a message: %s", d.Body)

			divisionTask := &devim_case.DivisionTask{}

			err := json.Unmarshal(d.Body, divisionTask)

			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}

		}
	}()

	// Stop for program termination
	<-stopChan
}