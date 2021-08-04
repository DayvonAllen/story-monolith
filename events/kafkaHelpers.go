package events

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/vmihailenco/msgpack/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"story-app-monolith/config"
	"story-app-monolith/domain"
)

func PushUserToQueue(message []byte, topic string) error {

	producer := GetInstance()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}


	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println(fmt.Errorf("%v", err))
		err = producer.Close()
		if err != nil {
			panic(err)
		}
		fmt.Println("Failed to send message to the queue")
	}

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", "user", partition, offset)
	return nil
}

func SendKafkaMessage(user *domain.User, eventType int) error {
	um := new(domain.Message)
	user.Password = ""
	um.User = *user

	// user created/updated event
	um.MessageType = eventType
	um.ResourceType = "user"

	fmt.Println(um.User)
	//turn user struct into a byte array
	b, err := msgpack.Marshal(um)

	if err != nil {
		return err
	}

	err = PushUserToQueue(b, config.Config("TOPIC"))

	if err != nil {
		return err
	}

	return nil
}

func HandleKafkaMessage(err error, user *domain.User, messageType int) error {
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return err
		}
		return fmt.Errorf("error processing data")
	}

	err = SendKafkaMessage(user, messageType)

	if err != nil {
		fmt.Println("Failed to publish new user")
	}

	return nil
}

func SendEventMessage(event *domain.Event, eventType int) error {
	um := new(domain.Message)
	um.Event = *event

	// user created/updated event
	um.MessageType = eventType
	um.ResourceType = "event"

	fmt.Println(um.Event)
	//turn user struct into a byte array
	b, err := msgpack.Marshal(um)

	if err != nil {
		return err
	}

	err = PushUserToQueue(b, "event")

	if err != nil {
		return err
	}

	return nil
}
