package shared

type configuration struct {
	AMQPConnectionURL string
}

var Config = configuration{
	AMQPConnectionURL: "amqp://admin:password@localhost:5672/",
}
