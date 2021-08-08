package repo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"story-app-monolith/database"
	"story-app-monolith/domain"
)

type NotificationRepoImpl struct {
	Notification             domain.Notification
	NotificationList             []domain.Notification
}

func (n NotificationRepoImpl) GetAllUnreadNotificationByUsername(username string) (*[]domain.Notification, error) {
	conn := database.MongoConn

	cur, err := conn.NotificationCollection.Find(context.TODO(), bson.D{{"for", username}, {"readStatus", false}})

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &n.NotificationList); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	err = cur.Close(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("error processing data")
	}


	return &n.NotificationList, nil
}


func NewNotificationRepoImpl() NotificationRepoImpl {
	var notificationRepoImpl NotificationRepoImpl

	return notificationRepoImpl
}
