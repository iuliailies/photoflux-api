package rabbitmq

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/iuliailies/photo-flux/internal/config"
	model "github.com/iuliailies/photo-flux/internal/models"
	"github.com/minio/minio-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

const QNAME = "upload"
const XNAME = "upload"
const VHOST = "photoflux"

type UploadsListener struct {
	db           *gorm.DB
	queueName    string
	exchangeName string
	connstr      string
}

func NewUploadsListener(db *gorm.DB, config config.RabbitMQ) UploadsListener {
	return UploadsListener{
		db:           db,
		queueName:    config.Queue,
		exchangeName: config.Exchange,
		connstr:      fmt.Sprintf("amqp://%s:%s@%s:%d/%s", config.User, config.Password, config.Host, config.Port, config.Vhost),
	}
}

// MinioNotification is a wrapper around the basic notification event, also
// containing the event name and the event key.
//
// This might not be necessary.
type MinioNotification struct {
	EventName string `json:"EventName"`
	Key       string `json:"Key"`
	minio.NotificationInfo
}

func (u UploadsListener) Start() error {

	// It's complicated to do the "initialization" step outside the thread
	// because there are multiple connections established for whom the
	// closing has to be deferred. But in the meantime, errors can happen
	// and this function has to be non-blocking.
	var errc = make(chan error)

	go func() {
		fmt.Println(u.connstr)
		conn, err := amqp.Dial(u.connstr)
		if err != nil {
			errc <- fmt.Errorf("could not connect to rabbitmq: %w", err)
			return
		}
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			errc <- fmt.Errorf("could not open channel: %w", err)
		}
		defer ch.Close()

		err = ch.ExchangeDeclare(
			u.exchangeName, // name
			"direct",       // type
			true,           // durable
			false,          // auto-deleted
			false,          // internal
			false,          // no-wait
			nil,            // arguments
		)
		if err != nil {
			errc <- fmt.Errorf("could not declare exchange")
			return
		}

		// We create a Queue to read messages from.
		q, err := ch.QueueDeclare(
			u.queueName, // name
			true,        // durable
			false,       // delete when unused
			false,       // exclusive
			false,       // no-wait
			nil,         // arguments
		)
		if err != nil {
			errc <- fmt.Errorf("could not declare queue: %w", err)
			return
		}

		// Tells RabbitMQ to not dispatch more than one message to a worker at a
		// time.
		err = ch.Qos(
			1,     // prefetch count
			0,     // prefetch size
			false, // global
		)
		if err != nil {
			errc <- fmt.Errorf("could not set up fair dispatch: %w", err)
			return
		}

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			false,  // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		if err != nil {
			errc <- fmt.Errorf("could not start consume messages: %w", err)
			return
		}

		// The main function can return while this goroutine is still running.
		close(errc)

		for d := range msgs {

			fmt.Printf("Received a message: %v\n", string(d.Body))
			// Acknowledge the message so RabbitMQ removes it from the queue.
			notif := MinioNotification{}
			err := json.Unmarshal(d.Body, &notif)
			if err != nil {
				// fmt.Println("could not marshal notification", err.Error())
				panic(err)
			}

			data := strings.SplitN(notif.Key, "/", 2)
			bucket := data[0]
			file := data[1]

			fmt.Println("bucket", bucket)
			fmt.Println("file", file)

			err = u.db.Model(&model.Photo{}).Where("id = ?", file).Update("is_uploaded", true).Error

			if err != nil {
				// Since no ack is send, the event will be requeued at some point.s
				fmt.Printf("an error occured during notification handling: %s", err.Error())
			} else {
				// Photo was updated succesfully
				d.Ack(false) // acknowlledged
			}
		}
	}()

	for err := range errc {
		if err != nil {
			return fmt.Errorf("an error occured when connecting to rabbitmq: %w", err)
		}
	}

	return nil
}
