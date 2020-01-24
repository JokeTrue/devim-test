package devim_case

type Configuration struct {
	AMQPConnectionURL string
}

type DivisionTask struct {
	Number int
}

var Config = Configuration{
	AMQPConnectionURL: "amqp://guest:guest@localhost:5672/",
}
