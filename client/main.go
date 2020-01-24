package main

import (
	devim_case "devim-case"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"time"
)

func main() {
	conn, err := amqp.Dial(devim_case.Config.AMQPConnectionURL)
	devim_case.HandleError(err, "Can't connect to AMQP")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	devim_case.HandleError(err, "Can't create a AMQP Channel")

	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("Division", true, false, false, false, nil)
	handleError(err, "Could not declare `Division` queue")

	rand.Seed(time.Now().UnixNano())

	ticker := time.NewTicker(3 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				divisionTask := devim_case.DivisionTask{Number: rand.Intn(10_000)}
				body, err := json.Marshal(divisionTask)
				if err != nil {
					devim_case.HandleError(err, "Error encoding JSON")
				}

				err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "text/plain",
					Body:         body,
				})

				if err != nil {
					log.Fatalf("Error publishing message: %s", err)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
