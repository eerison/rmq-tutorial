// https://www.rabbitmq.com/tutorials/tutorial-one-go.html
package main

import (
  "context"
  "log"
  "time"
  "os"
  "strings"

  amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
	  log.Panicf("%s: %s", msg, err)
	}
}

func bodyFrom(args []string) string {
    if (len(args) < 2) || os.Args[1] == "" {
        return "hello"
    }
    
	return strings.Join(args[1:], " ")
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
failOnError(err, "Failed to connect to RabbitMQ")
defer conn.Close()

ch, err := conn.Channel()
failOnError(err, "Failed to open a channel")
defer ch.Close()

q, err := ch.QueueDeclare(
	"hello", // name
	true,   // durable
	false,   // delete when unused
	false,   // exclusive
	false,   // no-wait
	nil,     // arguments
  )
  failOnError(err, "Failed to declare a queue")
  
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

body := bodyFrom(os.Args)
err = ch.PublishWithContext(ctx,
  "",           // exchange
  q.Name,       // routing key
  false,        // mandatory
  false,
  amqp.Publishing {
    DeliveryMode: amqp.Persistent,
    ContentType:  "text/plain",
    Body:         []byte(body),
  })
failOnError(err, "Failed to publish a message")
log.Printf(" [x] Sent %s", body)
}
