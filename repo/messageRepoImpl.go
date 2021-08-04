package repo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"story-app-monolith/database"
	"story-app-monolith/domain"
)

type MessageRepoImpl struct {
	Message             domain.Message
	MessageList             []domain.Message
	Conversation        domain.Conversation
	User				domain.User
}

func (m MessageRepoImpl) Create(message *domain.Message) (*domain.Conversation, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)


	err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"username", message.To}}).Decode(&m.User)

	if err != nil {
		return nil, err
	}

	for _, v := range m.User.BlockList {
		if v == message.From {
			return nil, fmt.Errorf("error")
		}

		if v == message.To {
			return nil, fmt.Errorf("error")
		}
	}

	for _, v := range m.User.BlockByList {
		if v == message.From {
			return nil, fmt.Errorf("error")
		}

		if v == message.To {
			return nil, fmt.Errorf("error")
		}
	}

	message.Id = primitive.NewObjectID()

	_, err = conn.MessageCollection.InsertOne(context.TODO(), &message)

	if err != nil {
		return nil, fmt.Errorf("error processing data")
	}

	err = conn.MessageCollection.FindOne(context.TODO(),
		bson.D{{"_id", message.Id}}).Decode(&m.Message)

	if err != nil {
		return nil, err
	}

	conversation, err := ConversationRepoImpl{}.FindByOwner(m.Message)

	if err != nil  {
		if err == mongo.ErrNoDocuments {
			err := ConversationRepoImpl{}.Create(m.Message)
			if err != nil {
				return nil, err
			}
			conversation, err = ConversationRepoImpl{}.FindConversation(m.Message.From, m.Message.To)
			if err != nil {
				return nil, err
			}
			return conversation, nil
		}
		return nil, err
	}
	conversation.Messages = append(conversation.Messages, m.Message)

	err = ConversationRepoImpl{}.UpdateConversation(*conversation, m.Message)

	if err != nil {
		return nil, err
	}

	conversation, err =  ConversationRepoImpl{}.FindConversation(m.Message.From, m.Message.To)

	return conversation, nil
}

func (m MessageRepoImpl) DeleteByID(owner string, id primitive.ObjectID) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	filter := bson.D{{"_id", id}}

	err := conn.MessageCollection.FindOne(context.TODO(), filter).Decode(&m.Message)

	if err != nil {
		return err
	}

	filter = bson.D{{"owner", owner}, {"messages._id", id}}

	err = conn.ConversationCollection.FindOne(context.TODO(), filter).Decode(&m.Conversation)

	if err != nil {
		return err
	}

	update := bson.M{"$pull": bson.M{"messages": m.Message}}

	_, err  = conn.ConversationCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	return nil
}

func exists(id primitive.ObjectID, ids []domain.DeleteMessage) bool {
	for _, v := range ids {
		if v.Id == id {
			return true
		}
	}

	return false
}

func (m MessageRepoImpl) DeleteAllByIDs(owner string, messageIds []domain.DeleteMessage) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	if len(messageIds) == 0 {
		return fmt.Errorf("bad values")
	}

	fmt.Println(messageIds)
	filter := bson.D{{"owner", owner}, {"messages._id", bson.M{"$in": bson.A{messageIds[0].Id}}}}

	err := conn.ConversationCollection.FindOne(context.TODO(), filter).Decode(&m.Conversation)

	if err != nil {
		return err
	}

	fmt.Println(m.Conversation)

	mes := make([]domain.Message, 0, len(m.Conversation.Messages))
	for _, v := range m.Conversation.Messages {
		if !exists(v.Id, messageIds) {
			mes = append(mes, v)
		}
	}

	m.Conversation.Messages = mes

	update := bson.M{"$set": bson.M{"messages": m.Conversation.Messages}}

	_, err = conn.ConversationCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	return nil
}

func NewMessageRepoImpl() MessageRepoImpl {
	var messageRepoImpl MessageRepoImpl

	return messageRepoImpl
}