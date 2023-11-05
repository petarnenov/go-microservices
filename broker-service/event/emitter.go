package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	conn *amqp.Connection
}

func (e *Emitter) setup() error {
	ch, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return declareExchange(ch)
}

func (e *Emitter) Push(event string, severity string) error {
	ch, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.Publish(
		"logs_topic", //exchange
		severity,     //routing key
		false,        //mandatory
		false,        //immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)

	return nil
}

func NewEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		conn: conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}