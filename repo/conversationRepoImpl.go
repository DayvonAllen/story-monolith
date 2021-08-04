package repo

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"log"
	"story-app-monolith/database"
	"story-app-monolith/domain"
	"time"
)

type ConversationRepoImpl struct {
	Message          domain.Message
	User			 domain.User
	Conversation     domain.Conversation
	Conversation2    domain.Conversation
	ConversationList []domain.Conversation
	ConversationPreview []domain.ConversationPreview
}

func (c ConversationRepoImpl) Create(message domain.Message) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"username", message.To}}).Decode(&c.User)

	if err != nil {
		return err
	}

	for _, v := range c.User.BlockList {
		if v == message.From {
			return fmt.Errorf("error")
		}

		if v == message.To {
			return fmt.Errorf("error")
		}
	}

	for _, v := range c.User.BlockByList {
		if v == message.From {
			return fmt.Errorf("error")
		}

		if v == message.To {
			return fmt.Errorf("error")
		}
	}

	c.Conversation.Id = primitive.NewObjectID()
	c.Conversation.CreatedAt = time.Now()
	c.Conversation.Owner = message.From
	c.Conversation.From = message.From
	c.Conversation.To = message.To
	c.Conversation.UnreadCount = 1
	c.Conversation.Messages = append(c.Conversation.Messages, message)
	c.Conversation.UpdatedAt = time.Now()

	c.Conversation2.Id = primitive.NewObjectID()
	c.Conversation2.CreatedAt = time.Now()
	c.Conversation2.Owner = message.To
	c.Conversation2.From = message.To
	c.Conversation2.To = message.From
	c.Conversation2.UnreadCount = 1
	c.Conversation2.Messages = append(c.Conversation2.Messages, message)
	c.Conversation2.UpdatedAt = time.Now()

	_, err = conn.ConversationCollection.InsertOne(context.TODO(), c.Conversation)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	_, err = conn.ConversationCollection.InsertOne(context.TODO(), c.Conversation2)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (c ConversationRepoImpl) FindByOwner(message domain.Message) (*domain.Conversation, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	fmt.Println(message)
	filter := bson.D{{"owner", message.From}}

	err := conn.ConversationCollection.FindOne(context.TODO(),
		filter).Decode(&c.Conversation)

	if err != nil {
		return nil, err
	}

	return &c.Conversation, nil
}

func (c ConversationRepoImpl) FindConversation(owner, to string) (*domain.Conversation, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	if owner == to {
		return nil, fmt.Errorf("bad request")
	}

	err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"username", to}}).Decode(&c.User)

	if err != nil {
		return nil, err
	}

	for _, v := range c.User.BlockList {
		if v == owner {
			return nil, fmt.Errorf("error")
		}

		if v == to {
			return nil, fmt.Errorf("error")
		}
	}

	for _, v := range c.User.BlockByList {
		if v == owner {
			return nil, fmt.Errorf("error")
		}

		if v == to {
			return nil, fmt.Errorf("error")
		}
	}

	filter := bson.M{
		"owner": owner,
		"$or": []interface{}{
			bson.M{"to": to},
			bson.M{"from": to},
		},
	}

	err = conn.ConversationCollection.FindOne(context.TODO(),
		filter).Decode(&c.Conversation)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.Conversation.Messages = make([]domain.Message,0, 0)
			return &c.Conversation, nil
		}
		return nil, err
	}

	if len(c.Conversation.Messages) == 0 {
		return &c.Conversation, nil
	}

	if c.Conversation.Messages[len(c.Conversation.Messages) - 1].Read != true {
		mes := make([]domain.Message, 0, len(c.Conversation.Messages))
		for _, v := range c.Conversation.Messages {
			v.Read = true

			mes = append(mes, v)
		}

		c.Conversation.UnreadCount = 0

		update := bson.D{{"$set", bson.D{{"messages", &mes}, {"unreadCount", c.Conversation.UnreadCount}}}}

		_, err = conn.ConversationCollection.UpdateMany(context.TODO(),
			filter, update)

		if err != nil {
			return nil, err
		}

		filter = bson.M{
			"owner": owner,
			"$or": []interface{}{
				bson.M{"to": to},
				bson.M{"from": to},
			},
		}

		err = conn.ConversationCollection.FindOne(context.TODO(),
			filter).Decode(&c.Conversation)
	}

	return &c.Conversation, nil
}

func (c ConversationRepoImpl) GetConversationPreviews(owner string) (*[]domain.ConversationPreview, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	filter := bson.M{
		"owner": owner,

	}

	cur, err := conn.ConversationCollection.Find(context.TODO(),
		filter)

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &c.ConversationList); err != nil {
		log.Fatal(err)
	}

	for _, v := range c.ConversationList {
		if len(v.Messages) > 0 {
			preview := new(domain.ConversationPreview)
			preview.Id = v.Id
			preview.To = v.To
			preview.PreviewMessage = v.Messages[len(v.Messages) - 1]
			preview.From = v.From
			preview.Owner = v.Owner
			preview.UnreadCount = v.UnreadCount
			preview.CreatedAt = v.CreatedAt
			preview.UpdatedAt = v.UpdatedAt
			c.ConversationPreview = append(c.ConversationPreview, *preview)
		}
	}

	return &c.ConversationPreview, nil
}


func (c ConversationRepoImpl) UpdateConversation(conversation domain.Conversation, message domain.Message) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.M{
		"owner": conversation.Owner,
		"$or": []interface{}{
			bson.M{"to": conversation.To},
			bson.M{"from": conversation.To},
		},
	}

	update := bson.D{{"$push", bson.D{{"messages", message}}},  {"$set", bson.D{{"updatedAt", time.Now()}}}}

	err := conn.ConversationCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&c.Conversation)

	if err != nil {
		return err
	}

	opts = options.FindOneAndUpdate().SetUpsert(true)
	filter = bson.M{
		"owner": conversation.To,
		"$or": []interface{}{
			bson.M{"to": conversation.Owner},
			bson.M{"from": conversation.Owner},
		},
	}
	update = bson.D{{"$push", bson.D{{"messages", message}}},  {"$set", bson.D{{"updatedAt", time.Now()}}},  {"$inc", bson.D{{"unreadCount", 1}}}}

	err = conn.ConversationCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&c.Conversation)

	if err != nil {
		return err
	}

	return nil
}

func (c ConversationRepoImpl) DeleteByID(conversationId primitive.ObjectID, username string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	res, err := conn.ConversationCollection.DeleteOne(context.TODO(), bson.D{{"_id", conversationId}, {"owner", username}})

	if err != nil {
		panic(err)
	}

	if res.DeletedCount == 0 {
		panic(fmt.Errorf("failed to delete story"))
	}

	return nil
}

func (c ConversationRepoImpl) DeleteAllByUsername(username string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	res, err := conn.ConversationCollection.DeleteMany(context.TODO(), bson.D{{"owner", username}})

	if err != nil {
		panic(err)
	}

	if res.DeletedCount == 0 {
		panic(fmt.Errorf("failed to delete story"))
	}

	return nil
}

func NewConversationRepoImpl() ConversationRepoImpl {
	var conversationRepoImpl ConversationRepoImpl

	return conversationRepoImpl
}
