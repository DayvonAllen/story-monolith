package repo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"story-app-monolith/database"
	"story-app-monolith/domain"
	"strconv"
	"sync"
	"time"
)

type ReadLaterRepoImpl struct {
	ReadLater     domain.ReadLater
	ReadLaterList []domain.ReadLater
	ReadLaterDto  domain.ReadLaterDto
}

func (r ReadLaterRepoImpl) Create(username string, storyId primitive.ObjectID) error {
	conn := database.MongoConn

	story := new(domain.StoryDto)
	err := conn.StoryCollection.FindOne(context.TODO(), bson.D{{"_id", storyId}}).Decode(&story)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return err
		}
		return fmt.Errorf("error processing data")
	}

	err = conn.ReadLaterCollection.FindOne(context.TODO(), bson.D{{"story._id", story.Id}, {"username", username}}).Decode(&story)

	if err != nil {
		r.ReadLater.CreatedAt = time.Now()
		r.ReadLater.UpdatedAt = time.Now()
		r.ReadLater.Username = username
		r.ReadLater.Story = *story
		r.ReadLater.Id = primitive.NewObjectID()

		_, err = conn.ReadLaterCollection.InsertOne(context.TODO(), &r.ReadLater)

		if err != nil {
			return fmt.Errorf("error processing data")
		}
		return nil
	}

	return fmt.Errorf("you already added this story to your read later list")
}

func (r ReadLaterRepoImpl) GetByUsername(username string, page string) (*domain.ReadLaterDto, error) {
	conn := database.MongoConn

	findOptions := options.FindOptions{}
	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}
	findOptions.SetSkip((int64(pageNumber) - 1) * int64(perPage))
	findOptions.SetLimit(int64(perPage))

	query := bson.D{{"username", username}}

	var wg sync.WaitGroup
	wg.Add(2)

	go func(conn *database.Connection) {
		defer wg.Done()
		cur, err := conn.ReadLaterCollection.Find(context.TODO(), query)

		if err != nil {
			panic(err)
		}

		if err = cur.All(context.TODO(), &r.ReadLaterList); err != nil {
			log.Fatal(err)
		}

		err = cur.Close(context.TODO())

		if err != nil {
			log.Fatal("error processing data")
		}

		return
	}(conn)

	go func(conn *database.Connection) {
		defer wg.Done()
		count, err := conn.ReadLaterCollection.CountDocuments(context.TODO(), query)

		if err != nil {
			panic(err)
		}

		r.ReadLaterDto.NumberOfStories = count

		if r.ReadLaterDto.NumberOfStories < 10 {
			r.ReadLaterDto.NumberOfPages = 1
		} else {
			r.ReadLaterDto.NumberOfPages = int(count/10) + 1
		}

		return
	}(conn)

	wg.Wait()

	r.ReadLaterDto.ReadLaterItems = r.ReadLaterList

	r.ReadLaterDto.CurrentPage = pageNumber

	return &r.ReadLaterDto, nil
}

func (r ReadLaterRepoImpl) Delete(id primitive.ObjectID, username string) error {
	conn := database.MongoConn

	res, err := conn.ReadLaterCollection.DeleteOne(context.TODO(), bson.D{{"_id", id}, {"username", username}})

	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("you can't delete this item")
	}

	return nil
}

func NewReadLaterRepoImpl() ReadLaterRepoImpl {
	var readLaterRepoImpl ReadLaterRepoImpl

	return readLaterRepoImpl
}
