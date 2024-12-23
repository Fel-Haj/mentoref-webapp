package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Account interface {
	GetID() primitive.ObjectID
	GetEmail() string
	GetPassword() string
}

type User struct {
	ID          primitive.ObjectID `bson:"_id"`
	Email       string             `bson:"email"`
	Password    string             `bson:"password"`
	FirstName   string             `bson:"first_name"`
	Surname     string             `bson:"surname"`
	Phone       string             `bson:"phone,omitempty"`
	CompanyName string             `bson:"company_name,omitempty"`
}

func (user *User) GetID() primitive.ObjectID {
	return user.ID
}

func (user *User) GetEmail() string {
	return user.Email
}

func (user *User) GetPassword() string {
	return user.Password
}

type Company struct {
	ID          primitive.ObjectID `bson:"_id"`
	Email       string             `bson:"email"`
	Password    string             `bson:"password"`
	CompanyName string             `bson:"company_name"`
}

func (company *Company) GetID() primitive.ObjectID {
	return company.ID
}

func (company *Company) GetEmail() string {
	return company.Email
}

func (company *Company) GetPassword() string {
	return company.Password
}

func ConnectDB(uri string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
