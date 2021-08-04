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
	"story-app-monolith/util"
	"strconv"
	"sync"
	"time"
)

type UserRepoImpl struct {
	users        []domain.User
	user         domain.User
	userDto      domain.UserDto
	userDtoList  []domain.UserDto
	userResponse domain.UserResponse
	currentUser domain.CurrentUserProfile
	viewedUser  domain.ViewUserProfile
}

func (u UserRepoImpl) FindAll(id primitive.ObjectID, page string, ctx context.Context, username string) (*domain.UserResponse, error) {
	var currentUser *domain.UserDto

	conn := database.MongoConn

	currentUser, err := u.FindByID(id, ctx)

	if err != nil {
		fmt.Println("Did not find in Cache in find all users...")
		return nil, err
	}

	findOptions := options.FindOptions{}
	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}
	findOptions.SetSkip((int64(pageNumber) - 1) * int64(perPage))
	findOptions.SetLimit(int64(perPage))

	// Get all users
	cur, err := conn.UserCollection.Find(ctx, bson.M{
		"profileIsViewable": true,
		"$and": []interface{}{
			bson.M{"_id": bson.M{"$ne": id}},
			bson.M{"_id": bson.M{"$nin": currentUser.BlockByList}},
			bson.M{"_id": bson.M{"$nin": currentUser.BlockList}},
		},
	}, &findOptions)

	if err != nil {
		return nil, err
	}

	if err = cur.All(ctx, &u.userDtoList); err != nil {
		log.Fatal(err)
	}

	u.userResponse = domain.UserResponse{Users: &u.userDtoList, CurrentPage: page}

	return &u.userResponse, nil
}

func (u UserRepoImpl) GetCurrentUserProfile(username string) (*domain.CurrentUserProfile, error) {
	conn := database.MongoConn

	err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&u.currentUser)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, fmt.Errorf("error processing data")
	}

	stories, err := StoryRepoImpl{}.FindAllByUsername(username)

	if err != nil {
		return nil, err
	}

	u.currentUser.Posts = *stories

	return &u.currentUser, nil
}

func (u UserRepoImpl) GetUserProfile(username, currentUsername string) (*domain.ViewUserProfile, error) {
	conn := database.MongoConn

	err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&u.viewedUser)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, fmt.Errorf("error processing data")
	}

	if u.viewedUser.ProfileIsViewable == false {
		return nil, fmt.Errorf("cannot view user")
	}

	if !u.viewedUser.DisplayFollowerCount {
		u.viewedUser.FollowerCount = -1
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		stories, err := StoryRepoImpl{}.FindAllByUsername(username)

		if err != nil {
			panic(err)
		}

		u.viewedUser.Posts = *stories

		return
	}()

	go func() {
		defer wg.Done()
		u.viewedUser.IsFollowing = helper.CurrentUserInteraction(u.viewedUser.Followers, currentUsername)
		return
	}()

	wg.Wait()

	return &u.viewedUser, nil
}

func (u UserRepoImpl) FindAllBlockedUsers(id primitive.ObjectID, ctx context.Context, username string) (*[]domain.UserDto, error) {

	var currentUser *domain.UserDto

	conn := database.MongoConn

	currentUser, err := u.FindByID(id, ctx)

	if err != nil {
		fmt.Println("Did not find in Cache in find all blocked users...")
		return nil, err
	}

	query := bson.M{"_id": bson.M{"$in": currentUser.BlockList}}

	// Get all users
	cur, err := conn.UserCollection.Find(context.TODO(), query)

	if err != nil {
		return nil, fmt.Errorf("error processing data")
	}

	var results []domain.UserDto
	if err = cur.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	u.userDtoList = results

	return &u.userDtoList, nil
}

func (u UserRepoImpl) Create(user *domain.User) error {
	conn := database.MongoConn

	cur, err := conn.UserCollection.Find(context.TODO(), bson.M{
		"$or": []interface{}{
			bson.M{"email": user.Email},
			bson.M{"username": user.Username},
		},
	})

	if err != nil {
		return fmt.Errorf("error processing data")
	}
	found := cur.Next(context.TODO())
	if !found {
		user.Id = primitive.NewObjectID()
		_, err = conn.UserCollection.InsertOne(context.TODO(), &user)

		if err != nil {
			return fmt.Errorf("error processing data")
		}

		return nil
	}
	err = cur.Decode(&u.userDto)
	if err != nil {
		return err
	}

	err = cur.Close(context.TODO())

	if err != nil {
		return err
	}

	if u.userDto.Username == user.Username {
		return fmt.Errorf("username is taken")
	}

	return fmt.Errorf("email is taken")
}

func (u UserRepoImpl) FindByID(id primitive.ObjectID, ctx context.Context) (*domain.UserDto, error) {
	conn := database.MongoConn

	err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&u.userDto)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, fmt.Errorf("error with the database")
	}

	return &u.userDto, nil
}

func (u UserRepoImpl) FindByUsername(username string, ctx context.Context) (*domain.UserDto, error) {
	conn := database.MongoConn

	err := conn.UserCollection.FindOne(context.TODO(), bson.M{"username": username, "$and":
	[]interface{}{
		bson.M{"profileIsViewable": true,
		},
	}}).Decode(&u.userDto)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("cannot find user")
		}
		return nil, fmt.Errorf("error processing data")
	}

	return &u.userDto, nil
}

func (u UserRepoImpl) UpdateByID(id primitive.ObjectID, user *domain.User) (*domain.UserDto, error) {
	conn := database.MongoConn

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"tokenHash", user.TokenHash}, {"tokenExpiresAt", user.TokenExpiresAt}}}}

	conn.UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts)

	return &u.userDto, nil
}

func (u UserRepoImpl) UpdateProfileVisibility(id primitive.ObjectID, user *domain.UpdateProfileVisibility, ctx context.Context) error {
	conn := database.MongoConn

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"profileIsViewable", user.ProfileIsViewable}}}}

	err := conn.UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&u.userDto)

	if err != nil {
		return err
	}

	u.userDto.ProfileIsViewable = user.ProfileIsViewable


	return nil
}

func (u UserRepoImpl) UpdateMessageAcceptance(id primitive.ObjectID, user *domain.UpdateMessageAcceptance, ctx context.Context) error {
	conn := database.MongoConn


	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"acceptMessages", user.AcceptMessages}}}}

	err := conn.UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&u.userDto)

	if err != nil {
		return err
	}

	u.userDto.AcceptMessages = user.AcceptMessages

	return nil
}

func (u UserRepoImpl) UpdateCurrentBadge(id primitive.ObjectID, user *domain.UpdateCurrentBadge, ctx context.Context) error {
	conn := database.MongoConn


	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"currentBadgeUrl", user.CurrentBadgeUrl}}}}

	err := conn.UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&u.userDto)

	if err != nil {
		return err
	}

	u.userDto.CurrentBadgeUrl = user.CurrentBadgeUrl

	return nil
}

func (u UserRepoImpl) UpdateProfilePicture(id primitive.ObjectID, user *domain.UpdateProfilePicture, ctx context.Context) error {
	conn := database.MongoConn


	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"profilePictureUrl", user.ProfilePictureUrl}}}}

	err := conn.UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&u.userDto)

	if err != nil {
		return err
	}

	u.userDto.ProfilePictureUrl = user.ProfilePictureUrl

	return nil
}

func (u UserRepoImpl) UpdateProfileBackgroundPicture(id primitive.ObjectID, user *domain.UpdateProfileBackgroundPicture, ctx context.Context) error {
	conn := database.MongoConn


	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"profileBackgroundPictureUrl", user.ProfileBackgroundPictureUrl}}}}

	err := conn.UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&u.userDto)

	if err != nil {
		return err
	}

	u.userDto.ProfileBackgroundPictureUrl = user.ProfileBackgroundPictureUrl

	return nil
}

func (u UserRepoImpl) UpdateCurrentTagline(id primitive.ObjectID, user *domain.UpdateCurrentTagline, ctx context.Context) error {
	conn := database.MongoConn

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"currentTagLine", user.CurrentTagLine}}}}

	err := conn.UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&u.userDto)

	if err != nil {
		return err
	}

	u.userDto.CurrentTagLine = user.CurrentTagLine

	return nil
}

func (u UserRepoImpl) UpdateDisplayFollowerCount(id primitive.ObjectID, user *domain.UpdateDisplayFollowerCount) error{
	conn := database.MongoConn

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"displayFollowerCount", user.DisplayFollowerCount}}}}

	err := conn.UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&u.userDto)

	if err != nil {
		return err
	}

	u.userDto.DisplayFollowerCount = user.DisplayFollowerCount

	return nil
}

func (u UserRepoImpl) UpdateVerification(id primitive.ObjectID, user *domain.UpdateVerification) error {
	conn := database.MongoConn

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"isVerified", user.IsVerified}}}}

	err := conn.UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&u.userDto)

	if err != nil {
		return err
	}

	u.userDto.IsVerified = user.IsVerified

	if err != nil {
		return err
	}

	return nil
}

func (u UserRepoImpl) UpdatePassword(id primitive.ObjectID, password string) error {
	conn := database.MongoConn

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"password", password}, {"tokenHash", ""}, {"tokenExpiresAt", 0}, {"updatedAt", time.Now()}}}}

	conn.UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts)

	return nil
}

func (u UserRepoImpl) UpdateFlagCount(flag *domain.Flag) error {
	conn := database.MongoConn

	cur, err := conn.FlagCollection.Find(context.TODO(), bson.M{
		"$and": []interface{}{
			bson.M{"flaggerID": flag.FlaggerID},
			bson.M{"flaggedUsername": flag.FlaggedUsername},
		},
	})

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	// todo send message
	if !cur.Next(context.TODO()) {
		flag.Id = primitive.NewObjectID()
		_, err = conn.FlagCollection.InsertOne(context.TODO(), &flag)

		if err != nil {
			return err
		}

		filter := bson.D{{"username", flag.FlaggedUsername}}
		update := bson.M{"$push": bson.M{"flagCount": flag.Id}}

		_, err = conn.UserCollection.UpdateOne(context.TODO(),
			filter, update)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("you've already flagged this user")
}

func (u UserRepoImpl) BlockUser(id primitive.ObjectID, username string, ctx context.Context, currentUsername string) error {
	conn := database.MongoConn

	err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&u.userDto)

	if id == u.userDto.Id {
		return fmt.Errorf("you can't block yourself")
	}

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("user not found")
		}
		return err
	}

	for _, foundUsername := range u.userDto.BlockByList {
		if foundUsername == username {
			return fmt.Errorf("already blocked")
		}
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

		// todo fix query
		filter := bson.D{{"_id", id}}
		update := bson.M{"$push": bson.M{"blockList": u.userDto.Username}}

		_, err = conn.UserCollection.UpdateOne(context.TODO(),
			filter, update)

		if err != nil {
			return nil, err
		}

		filter = bson.D{{"_id", u.userDto.Id}}
		update = bson.M{"$push": bson.M{"blockByList": currentUsername}}

		_, err = conn.UserCollection.UpdateOne(context.TODO(),
			filter, update)

		if err != nil {
			return nil, err
		}

		return nil, err
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)

	if err != nil {
		return fmt.Errorf("failed to block user")
	}

	go func() {
		user := new(domain.User)
		err = conn.UserCollection.FindOne(context.TODO(), bson.D{{"username", u.userDto.Username}}).Decode(user)

		if err != nil {
			panic(err)
			return
		}

		user2 := new(domain.User)

		err = conn.UserCollection.FindOne(context.TODO(), bson.D{{"username", currentUsername}}).Decode(user2)

		if err != nil {
			panic(err)
			return
		}

		return
	}()

	return nil
}

func (u UserRepoImpl) UnblockUser(id primitive.ObjectID, username string, ctx context.Context, currentUsername string) error {
	conn := database.MongoConn

	err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&u.userDto)

	if id == u.userDto.Id {
		return fmt.Errorf("you can't block or unblock yourself")
	}

	if err != nil {

		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("user not found")
		}
		return err
	}

	err = conn.UserCollection.FindOne(context.TODO(), bson.D{{"username", currentUsername}}).Decode(&u.user)

	if err != nil {
		return err
	}

	newBlockList, userIsBlocked := util.GenerateNewBlockList(username, u.user.BlockList)

	if !userIsBlocked {
		return fmt.Errorf("this user is not blocked")
	}

	currentUser := new(domain.UserDto)

	// todo better query
	err = conn.UserCollection.FindOne(context.TODO(), bson.D{{"username", currentUsername}}).Decode(&currentUser)

	blockList, userIsBlocked := util.GenerateNewBlockList(u.userDto.Username, currentUser.BlockList)

	if !userIsBlocked {
		return fmt.Errorf("this user is not blocked")
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

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {

		filter := bson.D{{"_id", id}}
		update := bson.M{"$set": bson.M{"blockList": blockList}}

		_, err = conn.UserCollection.UpdateOne(context.TODO(),
			filter, update)

		if err != nil {
			return nil, err
		}

		filter = bson.D{{"_id", u.userDto.Id}}
		update = bson.M{"$set": bson.M{"blockByList": newBlockList}}

		_, err = conn.UserCollection.UpdateOne(context.TODO(),
			filter, update)

		if err != nil {
			return nil, err
		}

		return nil, err
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)

	if err != nil {
		return fmt.Errorf("failed to unblock user")
	}

	go func() {
		user := new(domain.User)
		err = conn.UserCollection.FindOne(context.TODO(), bson.D{{"username", u.userDto.Username}}).Decode(user)

		if err != nil {
			panic(err)
			return
		}

		user2 := new(domain.User)

		err = conn.UserCollection.FindOne(context.TODO(), bson.D{{"username", currentUsername}}).Decode(user2)

		if err != nil {
			panic(err)
			return
		}

		return
	}()

	return nil
}

func (u UserRepoImpl) FollowUser(username string, currentUser string) error {
	conn := database.MongoConn

	err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"username", currentUser}}).Decode(&u.user)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	query := bson.M{"username": bson.M{"$in": u.user.Following}}

	// Get all users
	cur, err := conn.UserCollection.Find(context.TODO(), query)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	if cur.Next(context.TODO()){
		return fmt.Errorf("you are already following this user")
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

	var user = new(domain.User)
	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		filter := bson.D{{"username", currentUser}}
		update := bson.M{"$push": bson.M{"following": username}}

		_, err := conn.UserCollection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			return nil, err
		}

		err = conn.UserCollection.FindOne(context.TODO(), filter).Decode(&u.user)

		if err != nil {
			return nil, err
		}

		filter = bson.D{{"username", username}}
		update = bson.M{"$push": bson.M{"followers": currentUser}, "$inc": bson.M{"followerCount": 1}}

		_, err = conn.UserCollection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			return nil, err
		}

		err = conn.UserCollection.FindOne(context.TODO(), filter).Decode(&user)

		if err != nil {
			return nil, err
		}

		return nil, err
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)

	if err != nil {
		return err
	}

	return nil
}

func (u UserRepoImpl) UnfollowUser(username string, currentUser string) error {
	conn := database.MongoConn

	err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"username", currentUser}}).Decode(&u.user)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	query := bson.M{"username": bson.M{"$in": u.user.Following}}

	// Get all users
	cur, err := conn.UserCollection.Find(context.TODO(), query)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	if !cur.Next(context.TODO()){
		return fmt.Errorf("you are not following this user")
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

	var user = new(domain.User)
	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		filter := bson.D{{"username", currentUser}}
		update := bson.M{"$pull": bson.M{"following": username}}

		_, err := conn.UserCollection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			return nil, err
		}

		err = conn.UserCollection.FindOne(context.TODO(), filter).Decode(&u.user)

		if err != nil {
			return nil, err
		}

		filter = bson.D{{"username", username}}

		err = conn.UserCollection.FindOne(context.TODO(), filter).Decode(&user)

		if err != nil {
			return nil, err
		}

		update = bson.M{"$pull": bson.M{"followers": currentUser}, "$set": bson.M{"followerCount": user.FollowerCount - 1}}

		_, err = conn.UserCollection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			return nil, err
		}

		err = conn.UserCollection.FindOne(context.TODO(), filter).Decode(&user)

		if err != nil {
			return nil, err
		}

		return nil, err
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)

	if err != nil {
		return err
	}

	return nil
}

func (u UserRepoImpl) DeleteByID(id primitive.ObjectID, ctx context.Context, username string) error {
	conn := database.MongoConn

	_, err := conn.UserCollection.DeleteOne(context.TODO(), bson.D{{"_id", id}})

	if err != nil {
		return err
	}

	u.user.Id = id

	return nil
}

func NewUserRepoImpl() UserRepoImpl {
	var userRepoImpl UserRepoImpl

	return userRepoImpl
}
