package repo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"log"
	"story-app-monolith/database"
	"story-app-monolith/domain"
	helper "story-app-monolith/helpers"
	"sync"
	"time"
)

type ReplyRepoImpl struct {
	Reply        domain.Reply
	ReplyList    []domain.Reply
}

func (r ReplyRepoImpl) Create(comment *domain.Reply) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	commentObj := new(domain.Comment)

	err := conn.CommentsCollection.FindOne(context.TODO(), bson.D{{"_id", comment.ResourceId}}).Decode(&commentObj)

	if err != nil {
		return fmt.Errorf("resource not found")
	}

	_, err = conn.RepliesCollection.InsertOne(context.TODO(), &comment)

	if err != nil {
		return err
	}

	return nil
}

func (r ReplyRepoImpl) FindAllRepliesByResourceId(resourceID primitive.ObjectID, username string) (*[]domain.Reply, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	cur, err := conn.RepliesCollection.Find(context.TODO(), bson.D{{"resourceId", resourceID}})

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, fmt.Errorf("error processing data")
	}

	if err = cur.All(context.TODO(), &r.ReplyList); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	err = cur.Close(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("error processing data")
	}

	replies := make([]domain.Reply, 0, len(r.ReplyList))
	for _, v := range  r.ReplyList {
		v.CurrentUserLiked = helper.CurrentUserInteraction(v.Likes, username)
		if !v.CurrentUserLiked {
			v.CurrentUserDisLiked = helper.CurrentUserInteraction(v.Dislikes, username)
		}
		replies = append(replies, v)
	}
	return &replies, nil
}

func (r ReplyRepoImpl) UpdateById(id primitive.ObjectID, newContent string, edited bool, updatedTime time.Time, username string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}, {"authorUsername", username}}
	update := bson.D{{"$set", bson.D{{"content", newContent}, {"edited", edited},
		{"updatedTime", updatedTime}}}}

	err := conn.RepliesCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&r.Reply)

	if err != nil {
		return fmt.Errorf("cannot update comment that you didn't write")
	}

	return nil
}

func (r ReplyRepoImpl) LikeReplyById(commentId primitive.ObjectID, username string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	ctx := context.TODO()

	cur, err := conn.RepliesCollection.Find(ctx, bson.D{
		{"_id", commentId}, {"likes", username},
	})

	if err != nil {
		return err
	}

	if cur.Next(ctx) {
		return fmt.Errorf("you've already liked this comment")
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

		filter := bson.D{{"_id", commentId}}
		update := bson.M{"$pull": bson.M{"dislikes": username}}

		res, err := conn.RepliesCollection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			return nil, err
		}

		if res.MatchedCount == 0 {
			return nil, fmt.Errorf("cannot find story")
		}

		err = conn.RepliesCollection.FindOne(context.TODO(),
			filter).Decode(&r.Reply)

		r.Reply.DislikeCount = len(r.Reply.Dislikes)

		update = bson.M{"$push": bson.M{"likes": username}, "$inc": bson.M{"likeCount": 1}, "$set": bson.D{{"dislikeCount", r.Reply.DislikeCount}}}

		filter = bson.D{{"_id", commentId}}

		_, err = conn.RepliesCollection.UpdateOne(context.TODO(),
			filter, update)

		if err != nil {
			return nil, err
		}

		return nil, err
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)

	if err != nil {
		return fmt.Errorf("failed to like comment")
	}

	return nil
}

func (r ReplyRepoImpl) DisLikeReplyById(commentId primitive.ObjectID, username string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

	ctx := context.TODO()

	cur, err := conn.RepliesCollection.Find(ctx, bson.D{
		{"_id", commentId}, {"dislikes", username},
	})

	if err != nil {
		return err
	}

	if cur.Next(ctx) {
		return fmt.Errorf("you've already disliked this comment")
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

		filter := bson.D{{"_id", commentId}}
		update := bson.M{"$pull": bson.M{"likes": username}}

		res, err := conn.RepliesCollection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			return nil, err
		}

		if res.MatchedCount == 0 {
			return nil, fmt.Errorf("cannot find story")
		}

		err = conn.RepliesCollection.FindOne(context.TODO(),
			filter).Decode(&r.Reply)

		r.Reply.LikeCount = len(r.Reply.Likes)

		update = bson.M{"$push": bson.M{"dislikes": username}, "$inc": bson.M{"dislikeCount": 1}, "$set": bson.D{{"likeCount", r.Reply.LikeCount}}}

		filter = bson.D{{"_id", commentId}}

		_, err = conn.RepliesCollection.UpdateOne(context.TODO(),
			filter, update)

		if err != nil {
			return nil, err
		}

		return nil, err
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)

	if err != nil {
		return fmt.Errorf("failed to dislike comment")
	}

	return nil
}

func (r ReplyRepoImpl) UpdateFlagCount(flag *domain.Flag) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

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

	return fmt.Errorf("you've already flagged this comment")
}

func (r ReplyRepoImpl) DeleteById(id primitive.ObjectID, username string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)
	defer database.MongoConnectionPool.Put(conn)

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
		wg.Add(2)

		go func() {
			defer wg.Done()

			res, err := conn.RepliesCollection.DeleteOne(context.TODO(), bson.D{{"_id", id}, {"authorUsername", username}})

			if err != nil {
				panic(err)
			}

			if res.DeletedCount == 0  {
				panic(fmt.Errorf("failed to delete reply"))
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

		wg.Wait()

		return nil, err
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)

	if err != nil {
		return err
	}

	return nil
}

func NewReplyRepoImpl() ReplyRepoImpl {
	var replyRepoImpl ReplyRepoImpl

	return replyRepoImpl
}