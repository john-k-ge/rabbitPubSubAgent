package main

import (
	"log"
	"os"
)

func main() {
	var mode string

	args := os.Args

	switch len(args) {
	case 2:
		mode = args[1]
	default:
		mode = Sub
	}

	rabbit, err := NewRabbitAgent("amqp://0.0.0.0:5672", mode)
	if err != nil {
		log.Printf("Failed:  %v", err)
		return
	}

	switch mode {
	case Pub:
		errP := rabbit.Publish("HERP")
		if errP != nil {
			log.Printf("Failed to Pub: %v", err)
			return
		}
	case Sub:
		_, errS := rabbit.Subscribe()
		if errS != nil {
			log.Printf("Failed to Sub: %v", err)
			return
		}
	}

	//rabbitPub, err := NewRabbitAgent("amqp://0.0.0.0:5672", Pub)
	//
	//rabbitSub, err := NewRabbitAgent("amqp://0.0.0.0:5672", Sub)
	//if err != nil {
	//	log.Printf("Failed to create subscriber:  %v", err)
	//	return
	//}
	//errP := rabbitPub.Publish("HERP")
	//if errP != nil {
	//	log.Printf("Failed to Pub: %v", err)
	//	return
	//}
	//
	//_, errS := rabbitSub.Subscribe()
	//if errS != nil {
	//	log.Printf("Failed to Sub: %v", err)
	//	return
	//}

	log.Print("Done!")
}
