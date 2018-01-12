package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

const Pub = "publisher"
const Sub = "subscriber"
const exchange = "mypubsubexchange"

type RabbitPubSubAgent struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	mode       string
}

func NewRabbitAgent(amqp_url, mode string) (*RabbitPubSubAgent, error) {
	conn, err := amqp.Dial(amqp_url)
	if err != nil {
		log.Printf("Failed to dial AMQP: %v", err)
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to create channel: %v", err)
		return nil, err
	}

	if err := channel.ExchangeDeclare(exchange, "fanout", true, true, false, false, nil); err != nil {
		log.Printf("Failed to create fanout exchange `%v`: %v", exchange, err)
		return nil, err
	}

	rpsa := RabbitPubSubAgent{
		connection: conn,
		channel:    channel,
	}

	switch mode {
	case Pub:
		rpsa.mode = Pub
	default:
		rpsa.mode = Sub
	}

	return &rpsa, nil
}

func (p *RabbitPubSubAgent) Publish(message string) error {
	switch p.mode {
	case Sub:
		log.Printf("Cannot publish '%v%... I'm a subscriber!", message)
		return nil
	default:
		log.Printf("Going to publish %v", message)
		routingKey := "ignored for fanout exchanges, application dependent for other exchanges"
		err := p.channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
			Body: []byte(message),
		})
		p.channel.Close()
		if err != nil {
			log.Printf("Failed to publish message `%m`: %v", message, err)
			return err
		}
	}
	return nil
}

func (s *RabbitPubSubAgent) Subscribe() (string, error) {
	switch s.mode {
	case Pub:
		log.Print("I can't subscribe, I'm a publisher")
		return "", errors.New("I can't subscribe, I'm a publisher")
	default:
		host, _ := os.Hostname()
		if len(host) == 0 {
			host = "localhost"
		}

		myQueue := fmt.Sprintf("%v-%v", host, os.Getpid())
		if _, err := s.channel.QueueDeclare(myQueue, true, true, true, false, nil); err != nil {
			log.Printf("cannot consume from exclusive queue: %q, %v", myQueue, err)
			return "", err
		}

		routingKey := "application specific routing key for fancy toplogies"
		if err := s.channel.QueueBind(myQueue, routingKey, exchange, false, nil); err != nil {
			log.Printf("cannot consume without a binding to exchange: %q, %v", exchange, err)
			return "", err
		}

		deliveries, err := s.channel.Consume(myQueue, "", false, true, false, false, nil)
		if err != nil {
			log.Printf("cannot consume from: %q, %v", myQueue, err)
			return "", err
		}

		log.Printf("I subscribed to %v", myQueue)

		for msg := range deliveries {
			log.Printf("I received %s", msg.Body)
			return string(msg.Body), nil
		}
	}
	return "", nil
}
