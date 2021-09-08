package repo

import (
	"context"
	"story-app-monolith/database"
	"story-app-monolith/domain"
	helper "story-app-monolith/helpers"

	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"log"
	"strconv"
	"sync"
	"time"
)

type StoryRepoImpl struct {
	Story             domain.Story
	StoryDto          domain.StoryDto
	FeaturedStoryList []domain.FeaturedStoryDto
	StoryList         []domain.Story
	StoryDtoList      []domain.StoryDto
	StoryPreviews     []domain.StoryPreviewDto
	StoryPreviewList  domain.StoryList
}

func (s StoryRepoImpl) Create(story *domain.CreateStoryDto) error {
	conn := database.MongoConn

	story.Id = primitive.NewObjectID()

	if len(story.Content) > 160 {
		story.Preview = string([]rune(story.Content)[:161]) + "..."
	} else {
		story.Preview = story.Content
	}

	_, err := conn.StoryCollection.InsertOne(context.TODO(), &story)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (s StoryRepoImpl) UpdateById(id primitive.ObjectID, newContent string, newTitle string, username string, tag *domain.Tag, updated bool) error {
	conn := database.MongoConn

	filter := bson.D{{"_id", id}, {"authorUsername", username}}
	update := bson.D{{"$set",
		bson.D{{"content", newContent},
			{"title", newTitle},
			{"updatedAt", time.Now()},
			{"tag", tag},
			{"updated", updated},
		},
	}}

	_, err := conn.StoryCollection.UpdateOne(context.TODO(),
		filter, update)

	if err != nil {
		return fmt.Errorf("you can't update a story you didn't write")
	}

	story := new(domain.Story)

	story.Id = id
	story.Title = newTitle

	return nil
}

func (s StoryRepoImpl) FindAll(page string, newStoriesQuery bool) (*domain.StoryList, error) {
	conn := database.MongoConn

	findOptions := options.FindOptions{}
	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}
	findOptions.SetSkip((int64(pageNumber) - 1) * int64(perPage))
	findOptions.SetLimit(int64(perPage))

	if newStoriesQuery {
		findOptions.SetSort(bson.D{{"createdAt", -1}})
	}

	query := bson.M{}

	var wg sync.WaitGroup
	wg.Add(2)

	go func(conn *database.Connection) {
		defer wg.Done()
		cur, err := conn.StoryCollection.Find(context.TODO(), query, &findOptions)

		if err != nil {
			panic(err)
		}

		if err = cur.All(context.TODO(), &s.StoryPreviews); err != nil {
			log.Fatal(err)
		}
		return
	}(conn)

	go func(conn *database.Connection) {
		defer wg.Done()
		count, err := conn.StoryCollection.CountDocuments(context.TODO(), query)

		if err != nil {
			panic(err)
		}

		s.StoryPreviewList.NumberOfStories = count

		if s.StoryPreviewList.NumberOfStories < 10 {
			s.StoryPreviewList.NumberOfPages = 1
		} else {
			s.StoryPreviewList.NumberOfPages = int(count/10) + 1
		}
		return
	}(conn)

	wg.Wait()

	s.StoryPreviewList.Stories = s.StoryPreviews

	s.StoryPreviewList.CurrentPage = pageNumber

	cur, err := conn.StoryCollection.Find(context.TODO(), bson.M{}, &findOptions)

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &s.StoryList); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	err = cur.Close(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("error processing data")
	}

	return &s.StoryPreviewList, nil
}

func (s StoryRepoImpl) FindAllByUsername(username string) (*[]domain.StoryDto, error) {
	conn := database.MongoConn

	cur, err := conn.StoryCollection.Find(context.TODO(), bson.D{{"authorUsername", username}})

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &s.StoryDtoList); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	err = cur.Close(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("error processing data")
	}

	return &s.StoryDtoList, nil
}

func (s StoryRepoImpl) FeaturedStories() (*[]domain.FeaturedStoryDto, error) {
	conn := database.MongoConn

	findOptions := options.FindOptions{}

	findOptions.SetLimit(3)
	findOptions.SetSort(bson.D{{"score", -1}})

	cur, err := conn.StoryCollection.Find(context.TODO(), bson.M{}, &findOptions)

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &s.FeaturedStoryList); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	err = cur.Close(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("error processing data")
	}

	return &s.FeaturedStoryList, nil
}

func (s StoryRepoImpl) LikeStoryById(storyId primitive.ObjectID, username string) error {
	conn := database.MongoConn

	ctx := context.TODO()

	cur, err := conn.StoryCollection.Find(ctx, bson.D{
		{"_id", storyId}, {"likes", username},
	})

	if err != nil {
		return err
	}

	if cur.Next(ctx) {
		return fmt.Errorf("you've already liked this story")
	}

	// sets mongo's read and write concerns
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	// set up for a transaction
	session, err := conn.StartSession()

	if err != nil {
		panic(err)
	}

	defer session.EndSession(context.Background())

	// execute this code in a logical transaction
	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {

		filter := bson.D{{"_id", storyId}}
		update := bson.M{"$pull": bson.M{"dislikes": username}}

		res, err := conn.StoryCollection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			return nil, err
		}

		if res.MatchedCount == 0 {
			return nil, fmt.Errorf("cannot find story")
		}

		err = conn.StoryCollection.FindOne(context.TODO(),
			filter).Decode(&s.Story)

		s.Story.DislikeCount = len(s.Story.Dislikes)
		s.Story.Score++

		update = bson.M{"$push": bson.M{"likes": username}, "$inc": bson.M{"likeCount": 1}, "$set": bson.D{{"dislikeCount", s.Story.DislikeCount}, {"score", s.Story.Score}}}

		filter = bson.D{{"_id", storyId}}

		_, err = conn.StoryCollection.UpdateOne(context.TODO(),
			filter, update)

		if err != nil {
			return nil, err
		}

		return nil, err
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)

	if err != nil {
		return fmt.Errorf("failed to like story")
	}

	return nil
}

func (s StoryRepoImpl) DisLikeStoryById(storyId primitive.ObjectID, username string) error {
	conn := database.MongoConn

	ctx := context.TODO()

	cur, err := conn.StoryCollection.Find(ctx, bson.D{
		{"_id", storyId}, {"dislikes", username},
	})

	if err != nil {
		return err
	}

	if cur.Next(ctx) {
		return fmt.Errorf("you've already disliked this story")
	}

	// sets mongo's read and write concerns
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	// set up for a transaction
	session, err := conn.StartSession()

	if err != nil {
		panic(err)
	}

	defer session.EndSession(context.Background())

	// execute this code in a logical transaction
	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {

		filter := bson.D{{"_id", storyId}}
		update := bson.M{"$pull": bson.M{"likes": username}}

		res, err := conn.StoryCollection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			return nil, err
		}

		if res.MatchedCount == 0 {
			return nil, fmt.Errorf("cannot find story")
		}

		err = conn.StoryCollection.FindOne(context.TODO(),
			filter).Decode(&s.Story)

		s.Story.LikeCount = len(s.Story.Likes)
		s.Story.Score--

		update = bson.M{"$push": bson.M{"dislikes": username}, "$inc": bson.M{"dislikeCount": 1}, "$set": bson.D{{"likeCount", s.Story.LikeCount}, {"score", s.Story.Score}}}

		filter = bson.D{{"_id", storyId}}

		_, err = conn.StoryCollection.UpdateOne(context.TODO(),
			filter, update)

		if err != nil {
			return nil, err
		}

		return nil, err
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)

	if err != nil {
		return fmt.Errorf("failed to dislike story")
	}

	return nil
}

func (s StoryRepoImpl) FindById(storyID primitive.ObjectID, username string, userIp string) (*domain.StoryDto, error) {
	conn := database.MongoConn

	err := conn.StoryCollection.FindOne(context.TODO(), bson.D{{"_id", storyID}}).Decode(&s.StoryDto)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, fmt.Errorf("error processing data")
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		s.StoryDto.CurrentUserLiked = helper.CurrentUserInteraction(s.StoryDto.Likes, username)

		if !s.StoryDto.CurrentUserLiked {
			s.StoryDto.CurrentUserDisLiked = helper.CurrentUserInteraction(s.StoryDto.Dislikes, username)
		}
		return
	}()

	go func() {
		defer wg.Done()
		s.StoryDto.Comments, err = CommentRepoImpl{}.FindAllCommentsByResourceId(s.StoryDto.Id, username)
		if err != nil {
			panic(err)
		}
		return
	}()

	go func() {
		defer wg.Done()

		hasher := new(domain.Authentication)

		identity := new(domain.Identity)

		identityArr, err := hasher.SignToken([]byte(userIp))

		if err != nil {
			panic("Couldn't hash identity")
		}

		err = conn.IdentityCollection.FindOne(context.TODO(), bson.D{{"identifier", identityArr}, {"storyId", storyID}}).Decode(&identity)

		if err != nil {
			_, err = conn.StoryCollection.UpdateOne(context.TODO(), bson.D{{"_id", storyID}}, bson.M{"$inc": bson.M{"views": 1}})

			if err != nil {
				fmt.Println(fmt.Sprintf("%v", err))
			}

			val, err := hasher.SignToken([]byte(userIp))

			if err != nil {
				panic("Couldn't hash identity")
			}

			identity.Identifier = val
			identity.Id = primitive.NewObjectID()
			identity.StoryId = storyID

			_, err = conn.IdentityCollection.InsertOne(context.TODO(), identity)

			if err != nil {
				panic("Couldn't save identity")
			}
		}

		return
	}()

	wg.Wait()

	return &s.StoryDto, nil
}

func (s StoryRepoImpl) UpdateFlagCount(flag *domain.Flag) error {
	conn := database.MongoConn

	cur, err := conn.FlagCollection.Find(context.TODO(), bson.M{
		"$and": []interface{}{
			bson.M{"flaggerID": flag.FlaggerID},
			bson.M{"flaggedResource": flag.FlaggedResource},
		},
	})

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	if !cur.Next(context.TODO()) {
		flag.Id = primitive.NewObjectID()
		_, err = conn.FlagCollection.InsertOne(context.TODO(), &flag)

		return nil
	}

	return fmt.Errorf("you've already flagged this story")
}

func (s StoryRepoImpl) DeleteById(id primitive.ObjectID, username string) error {
	conn := database.MongoConn

	// sets mongo's read and write concerns
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	// set up for a transaction
	session, err := conn.StartSession()

	if err != nil {
		panic(err)
	}

	defer session.EndSession(context.Background())

	// execute this code in a logical transaction
	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		var wg sync.WaitGroup
		wg.Add(4)

		go func() {
			defer wg.Done()
			res, err := conn.StoryCollection.DeleteOne(context.TODO(), bson.D{{"_id", id}, {"authorUsername", username}})

			if err != nil {
				panic(err)
			}

			if res.DeletedCount == 0 {
				panic(fmt.Errorf("failed to delete story"))
			}
			return
		}()

		go func() {
			defer wg.Done()
			_, err = conn.FlagCollection.DeleteMany(context.TODO(), bson.D{{"flaggedResource", id}})

			if err != nil {
				panic(err)
			}
			return
		}()

		go func() {
			defer wg.Done()
			_, err = conn.ReadLaterCollection.DeleteMany(context.TODO(), bson.D{{"story._id", id}})

			if err != nil {
				panic(err)
			}
			return
		}()

		go func() {
			defer wg.Done()
			err = CommentRepoImpl{}.DeleteManyById(id, username)

			if err != nil {
				panic(err)
			}
			return
		}()

		wg.Wait()

		return nil, err
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)

	if err != nil {
		return err
	}

	story := new(domain.Story)
	story.Id = id

	return nil
}

func NewStoryRepoImpl() StoryRepoImpl {
	var storyRepoImpl StoryRepoImpl

	return storyRepoImpl
}
